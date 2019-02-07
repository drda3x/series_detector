package main

import "fmt"
import "os"
import "net/http"
import "io/ioutil"
import "encoding/xml"
import "strings"

type Item struct {
    XMLName xml.Name `xml:"item"`
    Title string `xml:"title"`
    Link string `xml:"link"`
    Category string `xml:"category"`
}

type Items struct {
    Items []Item `xml:"channel>item"` 
}

func get_xml(url string, i *Items) (error) {
    resp, err := http.Get(url)
    defer resp.Body.Close()

    if err != nil {
        return fmt.Errorf("Error: %v", err)
    }

    b_data, err := ioutil.ReadAll(resp.Body)

    if err != nil {
        return fmt.Errorf("Error: %v", err)
    }

    err = xml.Unmarshal(b_data, &i)
    
    return err
}

func get_torrent(url string, client *http.Client) ([]byte, error) {
    PostData := strings.NewReader("")
    req, err := http.NewRequest("GET", url, PostData)
    if err != nil {
        return nil, err
    }

    req.Header.Set("Cookie", "uid=524372; usess=a869b5eeff20d26736a9614964627baf")
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }

    b_data, err := ioutil.ReadAll(resp.Body)
    
    return b_data, err
}

func get_existed_torrents(folder string) []os.FileInfo {
    files, _ := ioutil.ReadDir("/tmp")
    return files
}

func main() {

    dat, _ := ioutil.ReadFile("/tmp/series_list")
    series := strings.Split(string(dat), "\n")

    if len(series) == 0 {
        return
    }

    url := "http://retre.org/rssdd.xml"

    //series := []string{"Звездный путь: Дискавери", "Орвилл", "Волшебники", "Хороший доктор", "Обратная сторона"}
    res := "1080p"

    var items Items
    err := get_xml(url, &items)

    if err != nil {
        panic(err)
        fmt.Println(err)
    } 
    client := &http.Client{}
    for _, i := range items.Items {
        
        in_series := false
        for _, s := range series[:len(series)-1] {
            in_series = strings.Contains(i.Title, s) && strings.Contains(i.Title, res)
        }
        
        downloaded := false
        for _, f := range get_existed_torrents("/tmp") {
            downloaded = strings.Contains(f.Name(), i.Title)
        }

        if !in_series || downloaded {
            continue
        } 

        data, err := get_torrent(i.Link, client)

        file, err := os.Create(fmt.Sprintf("/tmp/%s.torrent", i.Title))
        if err != nil {
            fmt.Println(err)
        }
        defer file.Close()

        file.Write(data)
        file.Sync()

        if err != nil {
            fmt.Println(err)
        } else {
            fmt.Println(fmt.Sprintf("FILE: %s. DL-success", i.Title))
        }
    }

}
