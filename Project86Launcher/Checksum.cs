using System;
using System.Diagnostics;
using System.IO;
using System.Linq;
using System.Security.Cryptography;
using System.Threading.Tasks;
using System.Windows;

namespace Project86Launcher;

public class Checksum
{
    public enum MethodType
    {
        Sequential, // downlaod files one by one if checksum mismatch
        Parallel // download files after checksum check is done
    }
    private string _root;
    private string _remoteChecksumPath;
    private string _remoteVersion;
    
    private ulong _totalSize = 0;

    public const string ServerURL = "https://";
    public string Signature => $"Checksum v1.0 | {_totalSize}";
    public string GatherPath => Path.Combine(_root, "Gather");
    
    public const string LineSeparator = "|";
    
    public delegate void ChecksumMismatchHandler(string path, ulong size);
    public delegate void DownloadProgressHandler(ulong downloaded, ulong total);
    
    private StreamWriter? _gatherWriter;
    
    private MethodType _methodType = MethodType.Parallel;
    
    public event DownloadProgressHandler DownloadProgress;

    public event EventHandler<(float current, float complete)> ChecksumProgress;

    public Checksum(string root, string remoteChecksumPath, string remoteVersion)
    {
        _root = root;
        _remoteChecksumPath = remoteChecksumPath;
        _remoteVersion = remoteVersion;
    }

    public Task CheckAsync(MethodType methodType = MethodType.Parallel)
    {
        return Task.Run(() => Check(methodType));
    }
    public void Check(MethodType methodType = MethodType.Parallel)
    {
        int linesCount = File.ReadAllLines(_remoteChecksumPath).Length;
        var streamReader = new StreamReader(_remoteChecksumPath);
        _methodType = methodType;
        
        ChecksumMismatchHandler handler = methodType switch
        {
            MethodType.Sequential => DownloadFile,
            MethodType.Parallel => GatherFiles,
            _ => throw new ArgumentOutOfRangeException(nameof(methodType), methodType, null)
        };
        if (methodType == MethodType.Parallel)
        {
            if (!StartGatherFiles())
            {
                // checksum process already done
                streamReader.Close();
                DownloadGatheredFiles();
                return;
            }
            else
                HeaderGatherFiles(_remoteVersion);
        }

        int count = 0;
        while (!streamReader.EndOfStream)
        {
            ChecksumProgress?.Invoke(this, (count, linesCount));
            var line = streamReader.ReadLine()!;
            var split = line.Split(LineSeparator);
            
            var locaPpath = split[0];
            var checksum = split[1];
            var size = ulong.Parse(split[2]);
            
            var path = Path.Combine(_root, locaPpath);
            
            
            if (!Compare(path, checksum))
            {
                handler(locaPpath, size);
            }
            count++;
        }
        ChecksumProgress?.Invoke(this, (count, linesCount));

        if (methodType == MethodType.Parallel)
        {
            EndGatherFiles();
            DownloadGatheredFiles();
        }
    }

    private bool StartGatherFiles()
    {
        if (!File.Exists(GatherPath))
        {

            _gatherWriter = new StreamWriter(GatherPath, false) ;
            return true;
        }

        var file = File.ReadLines(GatherPath);
        if (file.Count() < 2)
        {
            File.Delete(GatherPath);
            _gatherWriter = new StreamWriter(GatherPath);
            return true;
        }
        var header = file.First();
        if (header != _remoteVersion)
        {
            File.Delete(GatherPath);
            _gatherWriter = new StreamWriter(GatherPath);
            return true;
        }

        var footer = file.Last();
        var footerSplit = footer.Split(" | ");
        if (footerSplit[0] != Signature.Split(" | ")[0])
        {
            File.Delete(GatherPath);
            _gatherWriter = new StreamWriter(GatherPath);
            return true;
        }
        _totalSize = ulong.Parse(footerSplit[1]);

        return false; // good file version
    }
    
    private void EndGatherFiles()
    {
        _gatherWriter!.WriteLine(Signature);
        _gatherWriter.Close();
    }
    
    private void HeaderGatherFiles(string version)
    {
        _gatherWriter!.WriteLine(version);
    }
    private void GatherFiles(string path, ulong size)
    {
        _totalSize += size;
        _gatherWriter!.WriteLine(path + "|" + size);
    }
    
    public void DownloadGatheredFiles()
    {
        var streamReader = new StreamReader(GatherPath);
        streamReader.ReadLine(); // skip header
        
        
        ulong downloaded = 0;
        AwsAPI.WriteObjectProgressEvent += (sender, args) =>
        {
            DownloadProgress?.Invoke((ulong)args.TransferredBytes + downloaded , _totalSize);
        };
        while (!streamReader.EndOfStream)
        {
            var line = streamReader.ReadLine()!;
            if (line == Signature)
                break;
            var split = line.Split(LineSeparator);
            
            var locaPpath = split[0];
            var size = ulong.Parse(split[1]);
            
            DownloadProgress?.Invoke(downloaded, _totalSize);
            
            DownloadFile(locaPpath, size);
            downloaded += size;
            
            DownloadProgress?.Invoke(downloaded, _totalSize);
        }
    }

    public void DownloadFile(string path, ulong size)
    {

        Debug.WriteLine($"Downloading {path}");
        var folderName = $"Project86-v{_remoteVersion}/";
        var sanitizePath = path.Replace('\\', '/');
        var request = AwsAPI.DownloadObjectFromBucketAsync("project-86", $"{folderName}{sanitizePath}",
            _root, folderName);
        request.Wait();
        if (!request.Result)
            MessageBox.Show($"Failed to download {path} from S3.", "Download Error", MessageBoxButton.OK,
                MessageBoxImage.Error);
        Debug.WriteLine($"Downloaded {path} completed");
    }

    public bool Compare(string path, string checksum)
    {
        var localChecksum = GetChecksum(path);
        if (localChecksum != checksum)
        {
            Debug.WriteLine($"Checksum mismatch for {path}! Expected {checksum}, got {localChecksum}");
            return false;
        }
        else
        {
            Debug.WriteLine($"Checksum match for {path}!");
            return true;
        }
    }
    
    public static string GetChecksum(string path)
    {
        if (!File.Exists(path))
            return "";
        using var sha256 = SHA256.Create();
        
        using var stream = File.OpenRead(path);
        var hash = sha256.ComputeHash(stream);
        return BitConverter.ToString(hash).Replace("-", "").ToLowerInvariant();
    }
    
    /// <summary>
    /// Downloads the checksum file from S3
    /// </summary>
    /// <param name="path"></param>
    /// <param name="version"></param>
    /// <returns>The path of the download file</returns>
    public static string DownloadChecksum(string path, string version)
    {
        var folderName = $"Project86-v{version}/";
        var request = AwsAPI.DownloadObjectFromBucketAsync("project-86", $"{folderName}checksum.txt",
            path, folderName);
        request.Wait();
        if (!request.Result)
            MessageBox.Show($"Failed to download {path} from S3.", "Download Error", MessageBoxButton.OK,
                MessageBoxImage.Error);
        return Path.Combine(path, "checksum.txt");
    }

}