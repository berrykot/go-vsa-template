package database

type Client struct{} //своя реализация

func (client *Client) OpenConnection() error {
	return nil
}

func NewClient() (*Client, error) {
	client := &Client{}

	return client, nil
}
