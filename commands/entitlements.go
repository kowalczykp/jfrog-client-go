package commands

import (
    "strings"
    "strconv"
    "github.com/JFrogDev/bintray-cli-go/utils"
)

func ShowDownloadKeys(bintrayDetails *utils.BintrayDetails, org string) {
    path := getDownloadKeysPath(bintrayDetails, org)
    resp, body := utils.SendGet(path, nil, bintrayDetails.User, bintrayDetails.Key)
    if resp.StatusCode != 200 {
        utils.Exit(resp.Status + ". " + utils.ReadBintrayMessage(body))
    }
    println("Bintray response: " + resp.Status)
    println(string(body))
}

func CreateDownloadKey(flags *DownloadKeyFlags, org string) {
    data := buildDownloadKeyJson(flags, true)
    url := getDownloadKeysPath(flags.BintrayDetails, org)
    resp, body := utils.SendPost(url, []byte(data), flags.BintrayDetails.User, flags.BintrayDetails.Key)
    if resp.StatusCode != 201 {
        utils.Exit(resp.Status + ". " + utils.ReadBintrayMessage(body))
    }
    println("Bintray response: " + resp.Status)
    println(string(body))
}

func UpdateDownloadKey(flags *DownloadKeyFlags, org string) {
    data := buildDownloadKeyJson(flags, false)
    url := getDownloadKeysPath(flags.BintrayDetails, org)
    url += "/" + flags.Id
    resp, body := utils.SendPatch(url, []byte(data), flags.BintrayDetails.User, flags.BintrayDetails.Key)
    if resp.StatusCode != 200 {
        utils.Exit(resp.Status + ". " + utils.ReadBintrayMessage(body))
    }
    println("Bintray response: " + resp.Status)
    println(string(body))
}

func buildDownloadKeyJson(flags *DownloadKeyFlags, create bool) string {
    var existenceCheck string
    var whiteCidrs string
    var blackCidrs string
    if flags.ExistenceCheckUrl != "" {
        existenceCheck = "\"existence_check\": {" +
            "\"url\": \"" + flags.ExistenceCheckUrl + "\"," +
           "\"cache_for_secs\": \"" + strconv.Itoa(flags.ExistenceCheckCache) + "\"" +
        "}"
    }

    if flags.WhiteCidrs != "" {
        whiteCidrs = "\"white_cidrs\": " + fixCidr(flags.WhiteCidrs)
    }
    if flags.BlackCidrs != "" {
        blackCidrs = "\"black_cidrs\": " + fixCidr(flags.BlackCidrs)
    }

    data := "{"
    if create {
        data += "\"id\": \"" + flags.Id + "\","
    }
    data += "\"expiry\": \"" + flags.Expiry + "\""

    if existenceCheck != "" {
        data += "," + existenceCheck
    }
    if whiteCidrs != "" {
        data += "," + whiteCidrs
    }
    if blackCidrs != "" {
        data += "," + blackCidrs
    }
    data += "}"

    return data
}

func DeleteDownloadKey(flags *DownloadKeyFlags, org string) {
    url := getDownloadKeysPath(flags.BintrayDetails, org)
    url += "/" + flags.Id
    resp, body := utils.SendDelete(url, flags.BintrayDetails.User, flags.BintrayDetails.Key)
    if resp.StatusCode != 200 {
        utils.Exit(resp.Status + ". " + utils.ReadBintrayMessage(body))
    }
    println("Bintray response: " + resp.Status)
}

func fixCidr(cidr string) string {
    split := strings.Split(cidr, ",")
    if len(split) != 2 {
        utils.Exit("Invalid cidr format: " + cidr)
    }
    return "[\"" + split[0] + "\",\"" + split[1] + "\"]"
}

func getDownloadKeysPath(bintrayDetails *utils.BintrayDetails, org string) string {
    if org == "" {
        return bintrayDetails.ApiUrl + "users/" + bintrayDetails.User + "/download_keys"
    }
    return bintrayDetails.ApiUrl + "orgs/" + org + "/download_keys"
}

type DownloadKeyFlags struct {
    BintrayDetails *utils.BintrayDetails
    Id string
    Expiry string
    ExistenceCheckUrl string
    ExistenceCheckCache int
    WhiteCidrs string
    BlackCidrs string
}