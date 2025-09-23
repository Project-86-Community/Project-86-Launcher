package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/template"
)

type Version struct {
	Major   int
	Minor   int
	Patch   int
	Build   int
	Version string
}

func parseVersion(versionStr string, commitCount int) Version {
	cleanVersion := strings.TrimPrefix(versionStr, "v")
	parts := strings.Split(cleanVersion, ".")

	if len(parts) != 3 {
		fmt.Println("Invalid version format, expected v0.1.2.2")
		os.Exit(1)
	}

	zero, _ := strconv.Atoi(parts[0])
	major, _ := strconv.Atoi(parts[1])
	minorPatch, _ := strconv.Atoi(parts[2])

	return Version{
		Major:   zero,
		Minor:   major,
		Patch:   minorPatch,
		Build:   commitCount,
		Version: cleanVersion,
	}
}

const versionTemplate = `{
    "FixedFileInfo": {
        "FileVersion": {
            "Major": {{.Major}},
            "Minor": {{.Minor}},
            "Patch": {{.Patch}},
            "Build": {{.Build}}
        },
        "ProductVersion": {
            "Major": {{.Major}},
            "Minor": {{.Minor}},
            "Patch": {{.Patch}},
            "Build": {{.Build}}
        },
        "FileFlagsMask": "3f",
        "FileFlags ": "00",
        "FileOS": "040004",
        "FileType": "01",
        "FileSubType": "00"
    },
    "StringFileInfo": {
        "Comments": "",
        "CompanyName": "Project 86 Community",
        "FileDescription": "A Launcher developed for Project-86-Community-Game for managing game files.",
        "FileVersion": "{{.Version}}.{{.Build}}",
        "InternalName": "project86launcher.exe",
        "LegalCopyright": "Copyright © Project 86 Community",
        "LegalTrademarks": "",
        "OriginalFilename": "project86launcher.exe",ico
        "PrivateBuild": "",
        "ProductName": "Project 86 Launcher",
        "ProductVersion": "{{.Version}}",
        "SpecialBuild": ""
    },
    "VarFileInfo": {
        "Translation": {
            "LangID": "0409",
            "CharsetID": "04B0"
        }
    },
    "IconPath": "../../assets/images/icon.ico",
    "ManifestPath": "app.manifest"
}`

const manifestTemplate = `
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly xmlns="urn:schemas-microsoft-com:asm.v1" manifestVersion="1.0">
  <assemblyIdentity
    type="win32"
    name="Github.com.Project-86-Community.Project-86-Launcher"
    version="{{.Version}}.{{.Build}}"
    processorArchitecture="*"/>
  
  <!-- UAC Control (No Admin Required) -->
  <trustInfo xmlns="urn:schemas-microsoft-com:asm.v3">
    <security>
      <requestedPrivileges>
        <requestedExecutionLevel
          level="asInvoker"
          uiAccess="false"/>
      </requestedPrivileges>
    </security>
  </trustInfo>

  <!-- DPI Awareness (Prevents Blurry Text on 4K) -->
  <application xmlns="urn:schemas-microsoft-com:asm.v3">
    <windowsSettings>
      <dpiAwareness xmlns="http://schemas.microsoft.com/SMI/2016/WindowsSettings">PerMonitorV2</dpiAwareness>
      <dpiAware xmlns="http://schemas.microsoft.com/SMI/2005/WindowsSettings">true</dpiAware>
    </windowsSettings>
  </application>

  <!-- Windows 10/11 Compatibility -->
  <compatibility xmlns="urn:schemas-microsoft-com:compatibility.v1">
    <application>
      <supportedOS Id="{8e0f7a12-bfb3-4fe8-b9a5-48fd50a15a9a}"/> <!-- Windows 11 -->
      <supportedOS Id="{1f676c76-80e1-4239-95bb-83d0f6d0da78}"/> <!-- Windows 10 -->
    </application>
  </compatibility>
</assembly>`

func generateFile(tmplContent string, outputPath string, data Version) error {
	tmpl, err := template.New("file").Parse(tmplContent)
	if err != nil {
		return fmt.Errorf("error parsing template: %w", err)
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating %s: %w", outputPath, err)
	}
	defer func() {
		err := outputFile.Close()
		if err != nil {
			fmt.Println(fmt.Errorf("error closing %s: %w", outputPath, err))
		}
	}()

	if err := tmpl.Execute(outputFile, data); err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	return nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run scripts/version.go <version> <commit-count>")
		os.Exit(1)
	}

	versionStr := os.Args[1]
	commitCount, _ := strconv.Atoi(os.Args[2])
	version := parseVersion(versionStr, commitCount)

	// Generate versioninfo.json
	if err := generateFile(versionTemplate, "./cmd/p86l/versioninfo.json", version); err != nil {
		fmt.Printf("Error generating versioninfo.json: %v\n", err)
		os.Exit(1)
	}

	// Generate manifest file
	if err := generateFile(manifestTemplate, "./cmd/p86l/app.manifest", version); err != nil {
		fmt.Printf("Error generating manifest: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Successfully generated versioninfo.json and app.manifest")
}
