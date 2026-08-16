package livebooleanstatemachine

type Ticket struct {
	Held bool
	Sold bool
}

func Hold(ticket *Ticket) {
	ticket.Held = true
	ticket.Sold = false
}

func Sell(ticket *Ticket) {
	ticket.Held = false
	ticket.Sold = true
}
