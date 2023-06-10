using System;
using System.Diagnostics;
using System.IO;
using System.IO.Compression;
using System.Net;
using System.Net.Http;
using System.Text;
using System.Threading.Tasks;
using System.Windows;

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
        private Version LocalVersion
        {
            get => _localVersion;
            set
            {
                _localVersion = value;
                VersionText.Text = _localVersion.ToString();
            }
        }
        private Version _remoteVersion;
        
        private LauncherStatus _status;
        
        internal LauncherStatus Status
        {
            get => _status;
            set
            {
                _status = value;
                switch (_status)
                {
                    case LauncherStatus.Ready:
                        PlayButton.Content = "Play";
                        break;
                    case LauncherStatus.DownloadingUpdate:
                        PlayButton.Content = "Downloading Update...";
                        break;
                    case LauncherStatus.DownloadingGame:
                        PlayButton.Content = "Downloading Game...";
                        break;
                    case LauncherStatus.Failed:
                        PlayButton.Content = "Download Failed - Retry";
                        break;
                    default:
                        throw new ArgumentOutOfRangeException();
                }
            }
        }
        
        public MainWindow()
        {
            InitializeComponent();
            _rootPath = Directory.GetCurrentDirectory();
            _versionFile = Path.Combine(_rootPath, "version.txt");
            _gameZip = Path.Combine(_rootPath, "Build.zip");
            _gameExe = Path.Combine(_rootPath, "Build", "Luminosité Eternelle.exe");
            _gamePath = Path.Combine(_rootPath, "Build");
            Debug.WriteLine("Game exe path: " + _gameExe);
            File.WriteAllText(_versionFile, "1.0.0");

        }

        private void MainWindow_OnContentRendered(object? sender, EventArgs e)
        {
            CheckForUpdates();
        }

        private void PlayButton_OnClick(object sender, RoutedEventArgs e)
        {
            MessageBox.Show("This is a test message box.", "Test", MessageBoxButton.OK, MessageBoxImage.Information);
            Debug.WriteLine("File exists: " + File.Exists(_gameExe));
           
            if (File.Exists(_gameExe) && Status == LauncherStatus.Ready)
            {
                Debug.WriteLine("Starting game");
                ProcessStartInfo startInfo = new ProcessStartInfo(_gameExe);
                startInfo.WorkingDirectory = Path.Combine(_rootPath, "Build");
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
                try
                {
                    GitHubResponse? response = NetworkClient.GetAsync(GitHubAPIInfo.LatestReleaseURL).Result;
                    if (response == null)
                    {
                        Status = LauncherStatus.Failed;
                        MessageBox.Show("Failed to check for updates: Response was null.", "Error", MessageBoxButton.OK, MessageBoxImage.Error);
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
                    MessageBox.Show($"Failed to check for updates: {e.Message}", "Error", MessageBoxButton.OK, MessageBoxImage.Error);
                }
            }
            else // No version file, so we need to download the game
            {
                DownloadGame();
            }
        }

        private void DownloadGame()
        {
            Status = LauncherStatus.DownloadingGame;
            FetchingAndExtractingGameData();
        }

        /// <summary>
        /// A temporary function to download and extract the game data.
        /// Until a way to only download the changed files is found, this will be used.
        /// </summary>
        private void FetchingAndExtractingGameData()
        {
            HttpClient client = new HttpClient();
            client.GetAsync($"https://github.com/{GitHubAPIInfo.Owner}/{GitHubAPIInfo.Repo}/releases/latest/download/Build.zip").ContinueWith(task =>
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
            FetchingAndExtractingGameData();
        }
    }

    struct Version
    {
        internal static Version Zero = new Version(0, 0, 0);
        private short _major;
        private short _minor;
        private short _subMinor;
        
        internal Version(short major, short minor, short subMinor)
        {
            _major = major;
            _minor = minor;
            _subMinor = subMinor;
        }

        internal Version(string version, bool hasVersionPrefix = false)
        {
            if (hasVersionPrefix)
                version = version.Substring(1);
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
            return a._major == b._major && a._minor == b._minor && a._subMinor == b._subMinor;
        }

        public static bool operator !=(Version a, Version b)
        {
            return !(a == b);
        }

        public override string ToString()
        {
            return $"{_major}.{_minor}.{_subMinor}";
        }
    }

    internal struct GitHubAPIInfo
    {
        internal const string Repo = "Luminosite-Eternelle-public";
        internal const string Owner = "Taliayaya";
        
        // ReSharper disable once InconsistentNaming
        internal static string LatestReleaseURL => $"https://api.github.com/repos/{Owner}/{Repo}/releases/latest";
    }
}