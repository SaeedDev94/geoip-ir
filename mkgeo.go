package main

import (
  "bufio"
  "encoding/csv"
  "fmt"
  "io/ioutil"
  "net/http"
  "os"
  "strings"
  "time"
)

func MkGeo() {
  directory := "geolite2"
  today := time.Now().Format("2006-01-02")
  createDirectoryIfNotExists(directory)
  createCountryCsv(directory+"/GeoLite2-Country-Locations-ir.csv")
  ipv4List, ipv4Error := getFileContent(directory+"/"+today+".ipv4", "https://raw.githubusercontent.com/herrbischoff/country-ip-blocks/master/ipv4/ir.cidr")
  if ipv4Error != nil {
    fmt.Println("ipv4Error:", ipv4Error)
    return
  }
  ipv6List, ipv6Error := getFileContent(directory+"/"+today+".ipv6", "https://raw.githubusercontent.com/herrbischoff/country-ip-blocks/master/ipv6/ir.cidr")
  if ipv6Error != nil {
    fmt.Println("ipv6Error:", ipv6Error)
    return
  }
  createIpCsv(directory+"/GeoLite2-Country-Blocks-IPv4.csv", strings.TrimSpace(ipv4List))
  createIpCsv(directory+"/GeoLite2-Country-Blocks-IPv6.csv", strings.TrimSpace(ipv6List))
}

func createDirectoryIfNotExists(directory string) {
  // Check if the directory exists
  if _, err := os.Stat(directory); os.IsNotExist(err) {
    // Directory does not exist, create it
    err := os.Mkdir(directory, 0755)
    if err != nil {
      fmt.Println("Error creating directory:", err)
      return
    }
    fmt.Println("Directory created successfully:", directory)
  } else if err != nil {
    // Some other error occurred
    fmt.Println("Error checking directory:", err)
    return
  } else {
    // Directory already exists
    fmt.Println("Directory already exists:", directory)
  }
}

func getFileContent(fileName string, url string) (string, error) {
  // Check if the file exists
  _, err := os.Stat(fileName)
  if os.IsNotExist(err) {
    // File does not exist, make HTTP request to get content
    resp, err := http.Get(url)
    if err != nil {
      return "", err
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
      return "", err
    }
    content := string(body)
    writeToFile(fileName, content)
    return content, nil
  } else if err != nil {
    // Some other error occurred
    return "", err
  }
  // File exists, read its content
  fileContent, err := ioutil.ReadFile(fileName)
  if err != nil {
    return "", err
  }
  return string(fileContent), nil
}

func writeToFile(fileName, content string) error {
  file, err := os.Create(fileName)
  if err != nil {
    return err
  }
  defer file.Close()
  writer := bufio.NewWriter(file)
  _, err = writer.WriteString(content)
  if err != nil {
    return err
  }
  // Flush the writer to ensure all content is written to the file
  err = writer.Flush()
  if err != nil {
    return err
  }
  return nil
}

func createCountryCsv(fileName string) error {
  data := [][]string{
    {"geoname_id", "locale_code", "continent_code", "continent_name", "country_iso_code", "country_name", "is_in_european_union"},
    {"130758", "ir", "AS", "Asien", "IR", "Iran", "0"},
  }
  return createCSV(fileName, data)
}

func createIpCsv(fileName string, list string) error {
  data := [][]string{
    {"network", "geoname_id"},
  }
  lines := strings.Split(list, "\n")
  for _, line := range lines {
    data = append(data, []string{line, "130758"})
  }
  return createCSV(fileName, data)
}

func createCSV(fileName string, data [][]string) error {
  // Create or open the CSV file
  file, err := os.Create(fileName)
  if err != nil {
    return err
  }
  defer file.Close()
  // Create a new CSV writer
  writer := csv.NewWriter(file)
  defer writer.Flush()
  // Write data to CSV file
  for _, record := range data {
    if err := writer.Write(record); err != nil {
      return err
    }
  }
  return nil
}
