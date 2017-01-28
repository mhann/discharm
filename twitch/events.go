package twitch

/*
 * -------------------------
 * - Channel online events -
 * -------------------------
 */

type channelOnlineEventHandler func(*ChannelOnline)

type ChannelOnline struct {
	Channel twitchStreamChannel
	Name    string
}

func (eh channelOnlineEventHandler) Type() string {
	return channelOnlineEventType
}

func (eh channelOnlineEventHandler) New() interface{} {
	return &ChannelOnline{}
}

func (eh channelOnlineEventHandler) Handle(i interface{}) {
	if t, ok := i.(*ChannelOnline); ok {
		eh(t)
	}
}

/*
 * --------------------------
 * - Channel offline events -
 * --------------------------
 */

type channelOfflineEventHandler func(*ChannelOffline)

type ChannelOffline struct {
	Channel twitchStreamChannel
	Name    string
}

func (eh channelOfflineEventHandler) Type() string {
	return channelOfflineEventType
}

func (eh channelOfflineEventHandler) New() interface{} {
	return &ChannelOffline{}
}

func (eh channelOfflineEventHandler) Handle(i interface{}) {
	if t, ok := i.(*ChannelOffline); ok {
		eh(t)
	}
}

/*
 * -----------------------------------------------
 * - Textual representations of different events -
 * -----------------------------------------------
 */

const (
	channelOnlineEventType  = "CHANNEL_ONLINE"
	channelOfflineEventType = "CHANNEL_OFFLINE"
)
