package chatbotStructs

type primePair struct {
	factor    int
	remainder int
}

type InputsDeploy struct {
	ChatID string `json:"chatID"`
}

type RequestToGithubDeploy struct {
	Ref    string       `json:"ref"`
	Inputs InputsDeploy `json:"inputs"`
}

type requestToLoad struct {
	Url       string `json:"url"`
	Number    int    `json:"number"`
	Frequency int    `json:"frequency"`
	ChatID    int64  `json:"chat_id"`
}

type InputsReplicaCount struct {
	ChatID       string `json:"chatID"`
	ReplicaCount string `json:"number_of_replicas"`
	CustomUrl    string `json:"customURL"`
}

type RequestToGithubReplicaCount struct {
	Ref    string             `json:"ref"`
	Inputs InputsReplicaCount `json:"inputs"`
}

type MessageFromGitHub struct {
	ChatID string `json:"chat_id"`
}
