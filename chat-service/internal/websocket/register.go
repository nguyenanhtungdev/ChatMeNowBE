package websocket

func (h *Hub) RegisterClient(client *Client) {
	h.register <- client
}
