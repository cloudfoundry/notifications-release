package fakes

type Clients struct {
	store map[string]Client
}

func NewClients() *Clients {
	return &Clients{
		store: make(map[string]Client),
	}
}

func (c Clients) Add(client Client) {
	c.store[client.ID] = client
}

func (c Clients) Get(id string) (Client, bool) {
	client, ok := c.store[id]
	return client, ok
}

func (c *Clients) Clear() {
	c.store = make(map[string]Client)
}

func (c Clients) Delete(id string) bool {
	_, ok := c.store[id]
	delete(c.store, id)
	return ok
}
