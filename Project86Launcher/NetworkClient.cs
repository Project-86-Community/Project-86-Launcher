using System;
using System.Diagnostics;
using System.Net.Http;
using System.Net.Http.Headers;
using System.Threading.Tasks;
using Newtonsoft.Json;

namespace Project86Launcher;

public static class NetworkClient
{
    public static async Task<GitHubResponse?> GetAsync(string url)
    {
        
        Debug.WriteLine("starting get request");
        var client = new HttpClient();
        client.DefaultRequestHeaders.Accept.Add(new MediaTypeWithQualityHeaderValue("application/vnd.github+json"));
        client.DefaultRequestHeaders.UserAgent.Add(new ProductInfoHeaderValue("Project86Launcher", "1.0"));
        var task = client.GetAsync(url);
        Debug.WriteLine("waiting for response" + task.Status);
        var responseMessage = Task.Run(() => task).Result;
        Debug.WriteLine("Received response: " + responseMessage);

        var jsonResponse = await responseMessage.Content.ReadAsStringAsync();
        Debug.WriteLine("Received response: " + jsonResponse);
        
        return JsonConvert.DeserializeObject<GitHubResponse>(jsonResponse);
    }
    
    
    // Write a function that fetch the github latest release json of an url and deserialize it to a GitHubResponse object
    
    
    
    
}