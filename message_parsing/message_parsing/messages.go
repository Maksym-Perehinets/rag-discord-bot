package message_parsing

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"sync"
)

// messagesIteration gets messages from a channel in iterations.
// messageLimit allows defining the number of messages in a single iteration
// returns number of messages in a chanel
func messagesCount(session *discordgo.Session, ch chan int, wg *sync.WaitGroup, messageLimit int, channel string) {
	defer wg.Done()
	defer close(ch)
	lastMessage := ""

	for {
		messages, err := session.ChannelMessages(channel, messageLimit, lastMessage, "", "")
		if err != nil {
			log.Printf("Error getting messages for channel %s: %v", channel, err)
			return
		}
		ch <- len(messages)
		if len(messages) < messageLimit {
			break
		}
		lastMessage = messages[len(messages)-1].ID
	}

}

func messagesIteration(session *discordgo.Session, messageLimit int, channel string) <-chan []*discordgo.Message {
	out := make(chan []*discordgo.Message)
	go func() {
		defer close(out)
		lastMessage := ""

		for {
			newMessages, err := session.ChannelMessages(channel, messageLimit, lastMessage, "", "")
			if err != nil {
				log.Printf("Error getting messages for channel %s: %v", channel, err)
				return
			}
			out <- newMessages
			if len(newMessages) < messageLimit {
				break
			}
			lastMessage = newMessages[len(newMessages)-1].ID
		}
	}()
	return out
}

func MessageCount(session *discordgo.Session, channels []*discordgo.Channel) (messageCount map[string]int) {
	messageCount = make(map[string]int)
	messageLimit := 100
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, channel := range channels {
		ch := make(chan int)
		wg.Add(2)

		go messagesCount(session, ch, &wg, messageLimit, channel.ID)

		go func(chanelID string, ch chan int) {
			total := 0
			for count := range ch {
				total += count
			}
			mu.Lock()
			messageCount[chanelID] = total
			mu.Unlock()
			wg.Done()
		}(channel.ID, ch)
	}

	wg.Wait()
	return
}

// GetChannelMessages retrieves all messages from the specified channels.
// returns a map where the key is the channel ID and the value is a slice of messages.
func GetChannelMessages(session *discordgo.Session, channels []*discordgo.Channel) <-chan map[string][]*discordgo.Message {
	out := make(chan map[string][]*discordgo.Message)
	var wg sync.WaitGroup
	var mu sync.Mutex

	go func() {
		messagesMap := make(map[string][]*discordgo.Message)

		for _, channel := range channels {
			wg.Add(1)
			go func(channelID string) {
				in := messagesIteration(session, 100, channelID)
				var messages []*discordgo.Message
				for msgs := range in {
					messages = append(messages, msgs...)
				}
				mu.Lock()
				messagesMap[channelID] = messages
				mu.Unlock()
				wg.Done()
			}(channel.ID)
		}

		wg.Wait()
		out <- messagesMap
	}()

	return out
}
