package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)


// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	fmt.Print("Enter authorization code: ")
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	json.NewEncoder(f).Encode(token)
}

func main() {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	//config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
    config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := calendar.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	now := time.Now()
	t := now.Format(time.RFC3339)
    
    fmt.Println(string(t))
    
    calendarId := "primary"
    
	events, err := srv.Events.List(calendarId).ShowDeleted(false).
		SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}
	fmt.Println("Upcoming events:")
	if len(events.Items) == 0 {
		fmt.Println("No upcoming events found.")
	} else {
		for _, item := range events.Items {
			date := item.Start.DateTime
			if date == "" {
				date = item.Start.Date
			}
			
			fmt.Println(item.Id)
			fmt.Printf("%v (%v)\n", item.Summary, date)
			fmt.Printf("Место: %v\n", item.Location)
			fmt.Printf("Организатор: %v<%v>\n", item.Creator.DisplayName, item.Creator.Email)
			fmt.Printf("hangoutLink: %v\n", item.HangoutLink)
            
            for _, at := range item.Attendees {
                fmt.Printf("Attendees: %s<%v> - %v\n", at.DisplayName, at.Email, at.ResponseStatus)
                fmt.Println(at)
            }
            
            for _, a:= range item.Attachments {
                fmt.Printf("%s - %s\n", a.Title, a.FileUrl)
//                 fmt.Printf("attachments: %v - \n", _.Title, _.FileUrl)
            }
            
            fmt.Printf("--------------\n")
		}
	}
	
/*
     event := &calendar.Event{
        Summary: "Google I/O 2015",
        Location: "800 Howard St., San Francisco, CA 94103",
        Description: "A chance to hear more about Google's developer products.",
        Start: &calendar.EventDateTime{
            DateTime: now.Format(time.RFC3339),
            TimeZone: "America/Los_Angeles",
        },
        End: &calendar.EventDateTime{
            DateTime: now.Add(time.Minute * 10).Format(time.RFC3339),
            TimeZone: "America/Los_Angeles",
        },
        Attendees: []*calendar.EventAttendee{
            &calendar.EventAttendee{Email:"lpage@example.com", DisplayName: "Page", Id: "dfjdsnnfkdfjhdksf"},
            &calendar.EventAttendee{Email:"sbrin@example.com"},
        },
    }
	
	
    event, err = srv.Events.Insert(calendarId, event).Do()
    if err != nil {
        log.Fatalf("Unable to create event. %v\n", err)
    }
    fmt.Printf("Event created: %s\n", event.HtmlLink)

*/
    
    id := "ld8l71aabpshi02j1n97pc55j8_20181028T160000Z"
    
    my_event, err := srv.Events.Get(calendarId, id).Do()
    my_attendees := my_event.Attendees
    a_id := -1
    
    for i, item := range my_attendees {
        fmt.Println("Id:", item.Id)
        if item.Email == "lpage@example.com" {
            a_id = i
            break
        }
    }
    
    if a_id == -1 {
        log.Fatalf("Nothing to update", 1)
    }
    
    fmt.Println(my_attendees[a_id])
    my_a := calendar.EventAttendee{Email:"lpage@example.com", ResponseStatus: "needsAction"}
    my_attendees[a_id] = &my_a
    
    my_event.Attendees = my_attendees
    
    // с помощью Patch можно обновить только то что нужно
    my_event, err = srv.Events.Patch(calendarId, id, &calendar.Event{Attendees: my_attendees}).Do()
    
    fmt.Println(err)

}
