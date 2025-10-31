package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig/v2"
	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
)

var (
	settings             []azappconfig.Setting
	app                  *tview.Application
	header               *Header
	acDropdown           *tview.DropDown
	info                 *tview.Grid
	client               *azappconfig.Client
	cred                 *azidentity.DefaultAzureCredential
	keysManager          KeysManager
	valuesManager        *ValuesManager
	settingSearchManager SearchManager
	viewMode             ValueDisplayMode
)

const (
	VANITY_LOGO = `    ___   _______    __
   /   | / ____/ |  / /
  / /| |/ /    | | / / 
 / ___ / /___  | |/ /  
/_/  |_\____/  |___/   
Azure App Config Viewer
Version 0.1.0
`
)

type acvConfig struct {
	ConfigServers []string `yaml:"servers"`
}

func main() {
	configServers := []string{}
	if len(os.Args) > 1 {
		configServers = append(configServers, os.Args[1])
	}

	// Attempt to get config file
	home, err := os.UserHomeDir()
	if err == nil {
		acvConfigFile := path.Join(home, ".acv")
		_, err = os.Stat(acvConfigFile)
		if err == nil {
			// Found a config file.
			data, err := os.ReadFile(acvConfigFile)
			if err != nil {
				panic(fmt.Errorf("error reading confg file at %s", acvConfigFile))
			}

			var config acvConfig
			err = yaml.Unmarshal(data, &config)
			if err != nil {
				panic(errors.Wrapf(err, "error unmarshalling confg file at %s", acvConfigFile))
			}

			for _, configServer := range config.ConfigServers {
				if !slices.Contains(configServers, configServer) {
					configServers = append(configServers, configServer)
				}
			}
		}
	}

	if len(configServers) == 0 {
		fmt.Println("No app configurations to open, exiting")
		os.Exit(1)
	}

	cred, err = azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}

	// Top stuff
	header := NewHeader(configServers)

	// Setting Search Bar
	settingSearchManager = NewSearchManager(
		func(p tview.Primitive) {
			app.SetFocus(p)
		},
		func(s string) {
			findSettings(s)
			updateKeysList()
			app.SetFocus(getKeysManager().keys)
			settingSearchManager.setSearchType(NoSearch)
		},
	)

	// Navigable list of setting keys
	keysManager = NewKeysManager(func(s string) {
		revisions, err := getSettingRevisions(s, client)
		if err != nil {
			// chill for now
			return
		}

		if viewMode == Standard {
			getValuesManager().setPrimaryRevisions(s, revisions)
		} else {
			getValuesManager().setDiffRightRevisions(s, revisions)
		}

	})

	// Display of revision history and values
	valuesManager = NewValuesManager(
		func() {
			// Escaping out of the revisions dropdown, restore focus to the keys list
			app.SetFocus(getKeysManager().keys)
		},
		func(p tview.Primitive) {
			app.SetFocus(p)
		},
	)

	valuesManager.setRenderType(Plain)

	// Page layout
	pageGrid := tview.NewGrid().
		SetRows(8, 3, 0, 3).
		SetColumns(-3, -4).
		SetBorders(false).
		AddItem(header.GetPrimitive(), 0, 0, 1, 2, 0, 0, false).
		AddItem(keysManager.box, 1, 0, 2, 1, 0, 0, true).
		AddItem(settingSearchManager.searchBox, 3, 0, 1, 1, 0, 0, false).
		AddItem(valuesManager.GetPrimitive(), 1, 1, 3, 1, 0, 0, false)

	pageGrid.
		SetBorderStyle(tcell.Style{}.Bold(true)).SetBackgroundColor(tcell.ColorBlack)

	app = tview.NewApplication().SetRoot(pageGrid, true)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Block Ctrl-C to exit
		if event.Key() == tcell.KeyCtrlC {
			return nil
		}

		if settingSearchManager.searchType == StringSearch || valuesManager.valueSearchManager.searchType == StringSearch {
			// We are actively searching, don't steal the keystrokes
			return event
		}

		switch event.Rune() {
		case 'c':
			copyValue()
			return nil
		case 's':
			app.SetFocus(acDropdown)
			return nil
		case 'q':
			app.Stop()
			return nil
		case '/':
			// Search setting keys or setting value, depending which (if either) is focused
			if app.GetFocus() == keysManager.keys {
				settingSearchManager.setSearching(StringSearch)
				return nil
			} else if app.GetFocus() == valuesManager.valueTextView {
				valuesManager.updateValueBasedOnView()
				valuesManager.valueSearchManager.setSearching(StringSearch)
				return nil
			}

		case 'r':
			// Clear and reset, i.e. re-fetch from source
			// Don't clear out the values view if we're in Diff mode
			if viewMode != Diff {
				valuesManager.reset()
			}

			// But always do clear the search field
			settingSearchManager.Reset()

			fetchSettings("*")
			updateKeysList()
			app.SetFocus(getKeysManager().keys)
			return nil

		case 'j':
			// Toggle JSON rendering
			if valuesManager.renderType != Json {
				valuesManager.setRenderType(Json)
			} else {
				valuesManager.setRenderType(Plain)
			}
			valuesManager.updateValueBasedOnView()
			return nil

		case 'd':
			// Diff the current displayed value (if any) with another setting value
			if viewMode == Standard {
				// IMPORTANT! You can't enter diff mode if there is not a primary setting selected
				if len(valuesManager.primaryRevisionSelector.revisions) == 0 {
					break
				}

				setDisplayMode(Diff)
				// Focus the keys list, because pickikng a second setting to
				// diff is always what you want after entering diff mode
				app.SetFocus(keysManager.keys)
			} else {
				setDisplayMode(Standard)
			}

		}

		return event
	})

	// Ready to go, get secrets from the initial vault selection
	acDropdown.SetCurrentOption(0)
	fetchSettings("*")

	if err := app.Run(); err != nil {
		panic(err)
	}

}

func connect(serverUri string) {
	var err error
	// Establish a connection to the Key Vault client
	client, err = azappconfig.NewClient(serverUri, cred, nil)
	if err != nil {
		panic(err)
	}
}

// fetchSettings uses the server's filtering to fetch settings based on a filter string
func fetchSettings(keyFilter string) {
	settingsPager := client.NewListSettingsPager(
		azappconfig.SettingSelector{
			KeyFilter:   to.Ptr(keyFilter),
			LabelFilter: to.Ptr("*"),
			Fields:      azappconfig.AllSettingFields(),
		},
		nil,
	)

	settings = []azappconfig.Setting{}

	for settingsPager.More() {
		resp, err := settingsPager.NextPage(context.Background())
		if err != nil {
			panic(errors.Wrap(err, "failed to get paged settigns"))
		}

		settings = append(settings, resp.Settings...)
	}
}

// findSettings iterates through the currently fetched settings looking for key name
// substring matches for searchString, and sets the results to be the current fetched settings
func findSettings(searchString string) {

	newSettings := reduce(
		settings,
		func(s azappconfig.Setting) bool {
			res := strings.Contains(*s.Key, searchString)
			return res
		},
	)
	settings = newSettings
}

func updateKeysList() {
	keys := arraymap(settings, func(s azappconfig.Setting) string { return *s.Key })
	getKeysManager().updateKeys(keys)
}

func getKeysManager() KeysManager {
	return keysManager
}

func getValuesManager() *ValuesManager {
	return valuesManager
}

func getSettingRevisions(settingName string, client *azappconfig.Client) ([]azappconfig.Setting, error) {
	pager := client.NewListRevisionsPager(
		azappconfig.SettingSelector{
			KeyFilter:   to.Ptr(settingName),
			LabelFilter: to.Ptr("*"),
			Fields:      azappconfig.AllSettingFields(),
		},
		nil,
	)

	revisions := []azappconfig.Setting{}

	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		if err != nil {
			return nil, errors.Wrap(err, "failed to get paged secret versions")
		}

		revisions = append(revisions, resp.Settings...)
	}

	return revisions, nil
}

func copyValue() {
	clipboard.WriteAll(valuesManager.valueTextView.GetText(false))
}

func setDisplayMode(mode ValueDisplayMode) {
	viewMode = mode
	valuesManager.SetDisplayMode(mode)

	if mode == Diff {
		keysManager.SetTitle("Selecting For Diff Value (green)")
	} else {
		keysManager.SetTitle("")
	}
}
