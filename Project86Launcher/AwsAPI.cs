using System;
using System.Diagnostics;
using System.Threading;
using System.Threading.Tasks;
using Amazon;
using Amazon.Runtime;
using Amazon.Runtime.Endpoints;
using Amazon.S3;
using Amazon.S3.Model;

namespace Project86Launcher;

public static class AwsAPI
{
    
    const string BucketName = "project-86";
    private static IAmazonS3 _client;
    
    public static event EventHandler<WriteObjectProgressArgs> WriteObjectProgressEvent; 

    static AwsAPI()
    {
        var config = new AmazonS3Config
        {
            ServiceURL = "https://s3.leviia.com",
        };
        _client = new AmazonS3Client(Settings.PublicKey, Settings.PrivateKey ,config);
    }

    /// <summary>
    /// Shows how to download an object from an Amazon S3 bucket to the
    /// local computer.
    /// </summary>
    /// <param name="client">An initialized Amazon S3 client object.</param>
    /// <param name="bucketName">The name of the bucket where the object is
    /// currently stored.</param>
    /// <param name="objectName">The name of the object to download.</param>
    /// <param name="filePath">The path, including filename, where the
    /// downloaded object will be stored.</param>
    /// <returns>A boolean value indicating the success or failure of the
    /// download process.</returns>
    public static async Task<bool> DownloadObjectFromBucketAsync(
        IAmazonS3 client,
        string bucketName,
        string objectName,
        string filePath,
        string folderName)
    {
        // Create a GetObject request
        var request = new GetObjectRequest
        {
            BucketName = bucketName,
            Key = objectName,
        };
        // Issue request and remember to dispose of the response
        var downloadTask = client.GetObjectAsync(request);
        using GetObjectResponse response = await downloadTask;
        response.WriteObjectProgressEvent += (sender, args) => { WriteObjectProgressEvent?.Invoke(sender, args); };

        try
        {
            // Save object to local file
            var objectPath = objectName.Replace(folderName, "");
            await response.WriteResponseStreamToFileAsync($"{filePath}\\{objectPath}", false, CancellationToken.None);
            return response.HttpStatusCode == System.Net.HttpStatusCode.OK;
        }
        catch (AmazonS3Exception ex)
        {
            Console.WriteLine($"Error saving {objectName}: {ex.Message}");
            return false;
        }
        
        
        
    }
    
    public static async Task<bool> DownloadObjectFromBucketAsync(string bucketName, string objectName, string filePath, string folderName)
    {
        return await DownloadObjectFromBucketAsync(_client, bucketName, objectName, filePath, folderName);
    }

}