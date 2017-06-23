package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"gopkg.in/mgo.v2"
	"strconv"
	"time"
)

var realmsUS = map[string]string{
	"aggramar":         "en-us",
	"alleria":          "en-us",
	"archimonde":       "en-us",
	"argent-dawn":      "en-us",
	"arthas":           "en-us",
	"azgalor":          "en-us",
	"azjolnerub":       "en-us",
	"blackhand":        "en-us",
	"blackrock":        "en-us",
	"bleeding-hollow":  "en-us",
	"bloodhoof":        "en-us",
	"bloodscalp":       "en-us",
	"bonechewer":       "en-us",
	"boulderfist":      "en-us",
	"bronzebeard":      "en-us",
	"burning-blade":    "en-us",
	"burning-legion":   "en-us",
	"cenarion-circle":  "en-us",
	"cenarius":         "en-us",
	"crushridge":       "en-us",
	"daggerspine":      "en-us",
	"dalaran":          "en-us",
	"darkspear":        "en-us",
	"deathwing":        "en-us",
	"destromath":       "en-us",
	"dethecus":         "en-us",
	"doomhammer":       "en-us",
	"draenor":          "en-us",
	"dragonblight":     "en-us",
	"dragonmaw":        "en-us",
	"draka":            "en-us",
	"durotan":          "en-us",
	"earthen-ring":     "en-us",
	"elune":            "en-us",
	"eonar":            "en-us",
	"eredar":           "en-us",
	"feathermoon":      "en-us",
	"frostwolf":        "en-us",
	"garona":           "en-us",
	"gilneas":          "en-us",
	"gorefiend":        "en-us",
	"gorgonnash":       "en-us",
	"hellscream":       "en-us",
	"hyjal":            "en-us",
	"icecrown":         "en-us",
	"illidan":          "en-us",
	"kargath":          "en-us",
	"kelthuzad":        "en-us",
	"khadgar":          "en-us",
	"kiljaeden":        "en-us",
	"kilrogg":          "en-us",
	"lightbringer":     "en-us",
	"lightnings-blade": "en-us",
	"llane":            "en-us",
	"lothar":           "en-us",
	"magtheridon":      "en-us",
	"malganis":         "en-us",
	"malygos":          "en-us",
	"mannoroth":        "en-us",
	"medivh":           "en-us",
	"nathrezim":        "en-us",
	"nerzhul":          "en-us",
	"perenolde":        "en-us",
	"proudmoore":       "en-us",
	"sargeras":         "en-us",
	"shadowmoon":       "en-us",
	"shadowsong":       "en-us",
	"shattered-hand":   "en-us",
	"silver-hand":      "en-us",
	"silvermoon":       "en-us",
	"skullcrusher":     "en-us",
	"spinebreaker":     "en-us",
	"stonemaul":        "en-us",
	"stormrage":        "en-us",
	"stormreaver":      "en-us",
	"stormscale":       "en-us",
	"suramar":          "en-us",
	"terenas":          "en-us",
	"thunderhorn":      "en-us",
	"thunderlord":      "en-us",
	"tichondrius":      "en-us",
	"uldum":            "en-us",
	"uther":            "en-us",
	"warsong":          "en-us",
	"whisperwind":      "en-us",
	"windrunner":       "en-us",
	"zuljin":           "en-us",
}

var dungeons = [12]string{
	"black-rook-hold",
	"cathedral-of-eternal-night",
	"court-of-stars",
	"darkheart-thicket",
	"eye-of-azshara",
	"halls-of-valor",
	"maw-of-souls",
	"neltharions-lair",
	"return-to-karazhan-lower",
	"return-to-karazhan-upper",
	"the-arcway",
	"vault-of-the-wardens",
}

func GetPage(realm string, region string, dungeon string) *goquery.Document {
	fmt.Printf("GetPage :: start %v %v %v\n", realm, region, dungeon)
	url := fmt.Sprintf("https://worldofwarcraft.com/%v/game/pve/leaderboards/%v/%v", region, realm, dungeon)

	doc, err := goquery.NewDocument(url)
	if err != nil {
		fmt.Printf("error %v %v \n", url, err)
	}

	fmt.Printf("GetPage :: %v %v %v\n", realm, region, dungeon)
	return doc
}

// AffixParser parse the page to retrieve the week affix.
func AffixParser(doc *goquery.Document) []string {
	var affix []string
	doc.Find(".Media-text .font-semp-medium-white").Each(func(week_affix_index int, week_affix_row *goquery.Selection) {
		affix = append(affix, week_affix_row.Text())
	})
	return affix[:3]

}

func RunParser(doc *goquery.Document, Metadata Metadata, RunMetadata RunMetadata, MongoRunCollection *mgo.Collection) {

	doc.Find(".SortTable-body .SortTable-row").Each(func(run_index int, run_row *goquery.Selection) {
		fmt.Printf("RunParser :: %v :: start row\n", run_index)
		var run = Run{Metadata: Metadata, RunMetadata: RunMetadata}

		run_row.Find(".SortTable-col").Each(func(run_details_index int, run_details_row *goquery.Selection) {
			fmt.Printf("RunParser :: %v :: start details row %v\n", run_index, run_details_index)

			switch run_details_index {

			// case 0 = rank
			// Mythic level
			case 1:
				level, err := strconv.ParseInt(run_details_row.Text(), 10, 0)
				if err == nil {
					run.Level = level
				} else {
					run.Level = 0
				}
				fmt.Printf("RunParser :: %v :: row level :: %v\n", run_index, run.Level)

			// Time
			case 2:
				run.Time = run_details_row.Text()
				fmt.Printf("RunParser :: %v :: row time :: %v\n", run_index, run.Time)

			// Party
			case 3:
				run_details_row.Find(".List-item").Each(func(party_index int, party_row *goquery.Selection) {
					var player Player

					ref, _ := party_row.Find("a").Attr("href")
					role, _ := party_row.Find("a").Find(".Icon").Attr("class")

					// fishy way to determine the role of each people using the css class
					switch role {
					case "Icon Icon--role-tank Icon--small":
						player.Role = "tank"
					case "Icon Icon--role-healer Icon--small":
						player.Role = "healer"
					default:
						player.Role = "dps"
					}

					player.Armory = ref
					player.Name = party_row.Text()
					run.Party = append(run.Party, player)
					fmt.Printf("RunParser :: %v :: row party :: %v\n", run_index, player)
				})

			// Completed
			case 4:
				layout := "01/02/2006"
				t, _ := time.Parse(layout, run_details_row.Text())
				run.Completed = t
				fmt.Printf("RunParser :: %v :: row complet :: %v\n", run_index, run.Completed)
			}
		})
		run.Id = fmt.Sprintf("%v_%v_%v", run.Level, run.Time, run.Completed.Day())
		fmt.Printf("RunParser :: %v :: insert :: %v\n", run_index, run.Id)
		MongoRunCollection.Insert(run)
		fmt.Printf("RunParser :: %v :: insert done :: %v\n", run_index, run.Id)

	})
}

func GetMongoConnection() *mgo.Session {
	session, err := mgo.Dial("mongodb://localhost")

	if err != nil {
		panic(err)
	}

	return session
}

func main() {
	AffixPage := GetPage("aggramar", "en-us", "black-rook-hold")
	Affix := AffixParser(AffixPage)
	Metadata := Metadata{Patch: "7.2.5", Affix: Affix}

	MongoSession := GetMongoConnection()
	defer MongoSession.Close()
	MongoRunCollection := MongoSession.DB("test").C("run")

	for realm, region := range realmsUS {
		for _, dungeon := range dungeons {
			RunMetadata := RunMetadata{Realm: realm, Region: region, Dungeon: dungeon}
			doc := GetPage(realm, region, dungeon)
			RunParser(doc, Metadata, RunMetadata, MongoRunCollection)
		}
	}
}
