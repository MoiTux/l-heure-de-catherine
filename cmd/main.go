package main

import (
	"context"
	"fmt"

	"github.com/cloudevents/sdk-go/v2/event"
	l_heure_de_catherine "github.com/moitux/l-heure-de-catherine"
)

func main() {
	fmt.Println(l_heure_de_catherine.MidiSix(context.Background(), event.New()))
}
