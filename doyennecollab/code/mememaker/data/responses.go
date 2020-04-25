package data

//
// Response Data
//

// MemeGetResponse holds the data for a GET memes response from api.imgflip.com
type MemeGetResponse struct {
	Success bool     `json:"success"`
	Data    MemeData `json:"data"`
}

// MemeData defines the collection of memes
type MemeData struct {
	Memes []Meme `json:"memes"`
}

// Meme represents the details of an individual meme
type Meme struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	BoxCount int    `json:"box_count"`
}

// MemePostResponse handles the small response information
// returned from the POST call to api.imgflip.com
type MemePostResponse struct {
	Success bool        `json:"success"`
	Data    SuccessData `json:"data"`
	ErrMsg  string      `json:"error_message"`
}

// SuccessData holds teh links to the meme and the image after a successful call to
// POST api.imgflip.com
type SuccessData struct {
	URL     string `json:"url"`
	PageURL string `json:"page_url"`
}

// MemeError describes the error response received from api.imgflip.com
type MemeError struct {
	Success bool   `json:"success"`
	Data    string `json:"error_message"`
}
