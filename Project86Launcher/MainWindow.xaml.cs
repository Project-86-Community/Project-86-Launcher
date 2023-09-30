using System;
using System.ComponentModel;
using System.Diagnostics;
using System.IO;
using System.IO.Compression;
using System.Net;
using System.Net.Http;
using System.Text;
using System.Threading;
using System.Threading.Tasks;
using System.Windows;
using Amazon.S3;
using Amazon.S3.Model;
using Path = System.IO.Path;

namespace Project86Launcher
{
    enum LauncherStatus
    {
        Ready,
        DownloadingUpdate,
        DownloadingGame,
        Failed
    }
    /// <summary>
    /// Interaction logic for MainWindow.xaml
    /// </summary>
    public partial class MainWindow : Window
    {
        private string _rootPath;
        private string _gamePath;
        private string _versionFile;
        private string _gameZip;
        private string _gameExe;
        private Version _localVersion;
        
        private const string ExeName = "Project-86.exe";
        private Version LocalVersion
        {
            get => _localVersion;
            set
            {
                _localVersion = value;
                VersionText.Text = _localVersion.ToString();
                _gameZip = Path.Combine(_rootPath, $"Project86-v{_localVersion}.zip");
                _gameExe = Path.Combine(_rootPath, $"Build/", ExeName);
            }
        }
        private Version _remoteVersion;

        private Version FutureRemoteVersion
        {
            get => _remoteVersion;
            set
            {
                _gameZip = Path.Combine(_rootPath, $"Project86-v{_remoteVersion}.zip");
                _gameExe = Path.Combine(_rootPath, $"Build/", ExeName);
            }
        }
        
        private LauncherStatus _status;

        private string _buttonContent = "Play";

        internal LauncherStatus Status
        {
            get => _status;
            set
            {
                _status = value;
                switch (_status)
                {
                    case LauncherStatus.Ready:
                        _buttonContent = "Play";
                        break;
                    case LauncherStatus.DownloadingUpdate:
                        _buttonContent = "Downloading Update...";
                        break;
                    case LauncherStatus.DownloadingGame:
                        _buttonContent = "Downloading Game...";
                        break;
                    case LauncherStatus.Failed:
                        _buttonContent = "Download Failed - Retry";
                        break;
                    default:
                        throw new ArgumentOutOfRangeException();
                }
                Application.Current.Dispatcher.Invoke(() => PlayButton.Content = _buttonContent);
            }
        }
        public MainWindow()
        {
            InitializeComponent();
            _rootPath = Environment.GetFolderPath(Environment.SpecialFolder.ApplicationData);
            var path = Directory.CreateDirectory(Path.Combine(_rootPath, "Project-86-Community", "Project-86-Launcher"));
            _rootPath = path.FullName;
            _versionFile = Path.Combine(_rootPath, "version.txt");
            
            
            _gamePath = Path.Combine(_rootPath, $"Build/");
            //File.WriteAllText(_versionFile, "0.0.0-alpha"); // To force download last version
            if (!Directory.Exists(_gamePath))
                Directory.CreateDirectory(_gamePath);

        }

        private void MainWindow_OnContentRendered(object? sender, EventArgs e)
        {
            CheckForUpdates();
        }

        private void PlayButton_OnClick(object sender, RoutedEventArgs e)
        {
            Debug.WriteLine("File exists: " + File.Exists(_gameExe));
           
            if (File.Exists(_gameExe) && Status == LauncherStatus.Ready)
            {
                Debug.WriteLine("Starting game");
                ProcessStartInfo startInfo = new ProcessStartInfo(_gameExe);
                startInfo.WorkingDirectory = _gamePath;
                Process.Start(startInfo);
                
                Close();
            }
            else if (Status == LauncherStatus.Failed)
                CheckForUpdates();
        }

        private void CheckForUpdates()
        {
            if (File.Exists(_versionFile))
            {
                LocalVersion = new Version(File.ReadAllText(_versionFile));
            }
            else
                LocalVersion = Version.Zero;

            try
            {
                GitHubResponse? response = NetworkClient.GetAsync(GitHubAPIInfo.LatestReleaseURL).Result;
                if (response == null)
                {
                    Status = LauncherStatus.Failed;
                    MessageBox.Show("Failed to check for updates: Response was null.", "Error", MessageBoxButton.OK,
                        MessageBoxImage.Error);
                    return;
                }

                _remoteVersion = new Version(response.tag_name, true);
                Debug.WriteLine("Got latest release.");

                if (LocalVersion != _remoteVersion)
                {
                    DownloadUpdate();
                }
                else
                {
                    Status = LauncherStatus.Ready;
                }
            }
            catch (Exception e)
            {
                Status = LauncherStatus.Failed;
                MessageBox.Show($"Failed to check for updates: {e.Message}", "Error", MessageBoxButton.OK,
                    MessageBoxImage.Error);
            }
        }

        private void DownloadGame()
        {
            Status = LauncherStatus.DownloadingGame;
            //FetchingAndExtractingGameData();
            DownloadGameWithWebClient();
        }

        [Obsolete("Now using Amazon S3")]
        private void DownloadGameWithWebClient()
        {
            using (var client = new WebClient())
            {
                FutureRemoteVersion = _remoteVersion;
                client.DownloadFileAsync(new Uri(GitHubAPIInfo.LatestDownloadURL(_remoteVersion.ToString())), _gameZip);
                client.DownloadProgressChanged += ClientOnDownloadProgressChanged;
                client.DownloadFileCompleted += ClientOnDownloadFileCompleted;
            }
        }

        [Obsolete("Now using Amazon S3")]
        private void ClientOnDownloadFileCompleted(object? sender, AsyncCompletedEventArgs e)
        {
            LauncherStatus status;
            if (e.Error is not null)
            {
                Status = LauncherStatus.Failed;
                MessageBox.Show($"Failed to download game: {e.Error.Message}", "Error", MessageBoxButton.OK, MessageBoxImage.Error);
                return;
            }

            if (e.Cancelled)
            {
                Status = LauncherStatus.Failed;
                MessageBox.Show("Download cancelled.", "Error", MessageBoxButton.OK, MessageBoxImage.Error);
                return;
            }
            Status = LauncherStatus.Ready;
            ZipFile.ExtractToDirectory(_gameZip, _gamePath,  Encoding.UTF8, true); 
            File.Delete(_gameZip);
            File.WriteAllText(_versionFile, _remoteVersion.ToString());
            
            Application.Current.Dispatcher.Invoke(() => LocalVersion = _remoteVersion);
            
        }

        private void ClientOnDownloadProgressChanged(object sender, DownloadProgressChangedEventArgs e)
        {
            Application.Current.Dispatcher.Invoke(() => PlayButton.Content = _buttonContent + $" ({e.ProgressPercentage}%)");
        }
        
        private void OnDownloadProgress(ulong downloaded, ulong total)
        {
            var downloadMb = downloaded / 1024f / 1024f;
            var totalMb = total / 1024f / 1024f;
            
            Application.Current.Dispatcher.Invoke(() =>
            {
                DownloadProgressBar.Value = (downloadMb / totalMb) * 100;

                var downloadS = downloadMb > 1000 ? (downloadMb / 1024).ToString("F1") + "Gb" : downloadMb.ToString("F0") + "mb";
                var totalS = totalMb > 1000 ? (totalMb / 1024).ToString("F1") + "Gb" : totalMb.ToString("F0") + "mb";
                PlayButton.Content = _buttonContent +
                                            $" {downloadS}/{totalS} ({(int)(downloaded / (float)total * 100)}%)";
            });
        }

        private void OnChecksumProgress(object sender, (float current, float total) progress)
        {

            Application.Current.Dispatcher.Invoke(() =>
                {
                    PlayButton.Content = "Checking Files Integrity " + (progress.current / progress.total * 100).ToString("F0") + "/100%";
                    DownloadProgressBar.Value = progress.current / progress.total * 100;
                 
                });
        }
        

        /// <summary>
        /// A temporary function to download and extract the game data.
        /// Until a way to only download the changed files is found, this will be used.
        /// </summary>
        private void FetchingAndExtractingGameData()
        {
            HttpClient client = new HttpClient();
            
            client.GetAsync($"https://github.com/{GitHubAPIInfo.Owner}/{GitHubAPIInfo.Repo}/releases/latest/download/Project86-v.zip").ContinueWith(task =>
            {
                LauncherStatus status;
                if (task.IsFaulted)
                {
                    status = LauncherStatus.Failed;
                    MessageBox.Show($"Failed to download game: {task.Exception?.Message}", "Error", MessageBoxButton.OK, MessageBoxImage.Error);
                }
                else
                {
                    File.WriteAllBytes(_gameZip, task.Result.Content.ReadAsByteArrayAsync().Result);
                    status = LauncherStatus.Ready;
                    ZipFile.ExtractToDirectory(_gameZip, _gamePath,  Encoding.UTF8, true); 
                    File.Delete(_gameZip);
                    File.WriteAllText(_versionFile, _remoteVersion.ToString());
                    Application.Current.Dispatcher.Invoke(() => LocalVersion = _remoteVersion);
                }

                Application.Current.Dispatcher.Invoke(() => Status = status);
            });
        }
        
        
        




        private void DownloadUpdate()
        {
            
            Status = LauncherStatus.DownloadingUpdate;
            Debug.WriteLine("Downloading checksum file.");
            var checksumPath = Task.Run(() => Checksum.DownloadChecksum(_rootPath, _remoteVersion.ToString()));
            checksumPath.ContinueWith((task) =>
            {

                Debug.WriteLine("Checksum file downloaded.");
                // @"F:\CheckSumCreator\ChecksumCreator\ChecksumCreator\bin\Debug\net7.0\checksum.txt"
                Checksum checksum = new Checksum(_gamePath, task.Result, _remoteVersion.ToString());
                checksum.DownloadProgress += OnDownloadProgress;
                checksum.ChecksumProgress += OnChecksumProgress;
                checksum.CheckAsync().ContinueWith(
                    delegate
                    {
                        Debug.WriteLine("Checksum done.");
                        Status = LauncherStatus.Ready;
                        checksum.DownloadProgress -= OnDownloadProgress;
                        File.WriteAllText(_versionFile, _remoteVersion.ToString());
                        Application.Current.Dispatcher.Invoke(() => LocalVersion = _remoteVersion);
                    });

                //DownloadGameWithWebClient();
                //FetchingAndExtractingGameData();
            });
        }
    }

    struct Version
    {
        internal static Version Zero = new Version(0, 0, 0);
        private short _major;
        private short _minor;
        private short _subMinor;
        public bool IsAlpha;
        
        internal Version(short major, short minor, short subMinor, bool isAlpha = false)
        {
            _major = major;
            _minor = minor;
            _subMinor = subMinor;
            IsAlpha = isAlpha;
        }

        internal Version(string version, bool hasVersionPrefix = false)
        {
            if (hasVersionPrefix)
                version = version.Substring(1);
            if (version.Contains("-"))
            {
                IsAlpha = true;
                version = version.Substring(0, version.IndexOf('-'));
            }
            else
            {
                IsAlpha = false;
            }
            Debug.WriteLine("Version is " + version);
            string[] versionSplit = version.Split('.');
            if (versionSplit.Length != 3)
            {
                throw new ArgumentException("Version string must be in the format \"major.minor.subMinor\"");
            }

            _major = short.Parse(versionSplit[0]);
            _minor = short.Parse(versionSplit[1]);
            _subMinor = short.Parse(versionSplit[2]);
        }
        
        public static bool operator ==(Version a, Version b)
        {
            return a._major == b._major && a._minor == b._minor && a._subMinor == b._subMinor && a.IsAlpha == b.IsAlpha;
        }

        public static bool operator !=(Version a, Version b)
        {
            return !(a == b);
        }

        public override string ToString()
        {
            return $"{_major}.{_minor}.{_subMinor}" + (IsAlpha ? "-alpha" : "");
        }
    }

    internal struct GitHubAPIInfo
    {
        internal const string Repo = "Project-86";
        internal const string Owner = "Taliayaya";
        
        // ReSharper disable once InconsistentNaming
        internal static string LatestReleaseURL => $"https://api.github.com/repos/{Owner}/{Repo}/releases/latest";

        internal static string LatestDownloadURL(string version) =>
            $"https://github.com/{Owner}/{Repo}/releases/download/v{version}/Project86-v{version}.zip";
    }
}