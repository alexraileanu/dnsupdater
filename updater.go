package main

import (
    "./models"
    "context"
    "github.com/digitalocean/godo"
    "github.com/joho/godotenv"
    "io/ioutil"
    "net/http"
    "os"
    "strconv"
    "time"
)

func timestamp() string {
    now := time.Now()

    return now.Format("2006-01-02 15:04")
}

func logMessage(message string) string {
    return "[" + timestamp() + "] - " + message + "\n"
}

func getNewIP() string {
    response, _ := http.Get("https://api.ipify.org")
    responseBody, _ := ioutil.ReadAll(response.Body)

    return string(responseBody)
}

func checkIP() string {
    oldIP := ""

    if _, err := os.Stat("ip"); err == nil {
        fileContents, _ := ioutil.ReadFile("ip")
        oldIP = string(fileContents)
    }

    newIP := getNewIP()

    if oldIP != newIP {
        _ = ioutil.WriteFile("ip", []byte(newIP), 644)

        return newIP
    }

    return oldIP
}

func updateDNS(ip string) error {
    client := models.MakeClient()
    recordID, _ := strconv.ParseInt(os.Getenv("RECORD_ID"), 10, 32)

    editRecord := &godo.DomainRecordEditRequest{
        Data: ip,
    }

    _, _, err := client.Domains.EditRecord(context.Background(), os.Getenv("DOMAIN"), int(recordID), editRecord)

    return err
}

func main() {
    err := godotenv.Load()
    if err != nil {
        panic(err)
    }

    ip := checkIP()
    if ip != "" {
        print(logMessage("Attempting to update DNS record"))

        err := updateDNS(ip)

        if err != nil {
            panic(err)
        }

        print(logMessage("Updated DNS record to " + ip))
    } else {
        print(logMessage("Skipping request as IP has not changed"))
    }
}
