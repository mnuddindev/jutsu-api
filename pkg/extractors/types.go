package extractors

type Role struct {
	Anime struct {
		ID     string `json:"id"`
		Title  string `json:"title"`
		Poster string `json:"poster"`
		Type   string `json:"type"`
		Year   string `json:"year"`
	} `json:"anime"`
	Character struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Profile string `json:"profile"`
		Role    string `json:"role"`
	} `json:"character"`
}

type VoiceActorData struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Profile      string `json:"profile"`
	JapaneseName string `json:"japaneseName"`
	About        struct {
		Description string `json:"description"`
		Style       string `json:"style"`
	} `json:"about"`
	Roles []Role `json:"roles"`
}

type VoiceActorResult struct {
	Success bool `json:"success"`
	Results struct {
		Data []VoiceActorData `json:"data"`
	} `json:"results"`
}

type CharactersVoiceActors struct {
	Character struct {
		ID     string `json:"id"`
		Poster string `json:"poster"`
		Name   string `json:"name"`
		Cast   string `json:"cast"`
	} `json:"character"`
	VoiceActors []struct {
		ID     string `json:"id"`
		Poster string `json:"poster"`
		Name   string `json:"name"`
	} `json:"voiceActors"`
}
