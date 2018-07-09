package patent


type PatentResults struct {
    Title   string
    Is_application bool
    Snippet string
    Priority_date string
    Filing_date string
    Publication_date string
    Inventor string
    Assignee string 
    Publication_number string
    Language string
    Thumbnail string
    Pdf string
}


type PatentDetail struct {
	Prefix string
    ApplicationNumber   string
    Title string
    AssigneeOriginal string
    Language string
    PriorityDate string
    FilingDate string
    PublicationDate string
    GrantDate string
    Inventor string
    Assignee string 
    Pdf string
    RelevantPatents []PatentDetail
}

func (patent PatentDetail) String() string {
	return "" + patent.Prefix + ", " +
			patent.ApplicationNumber + ", " +
			patent.Title + ", " +
			patent.AssigneeOriginal + ", " +
			patent.Language + ", " +
			patent.PriorityDate + ", " +
			patent.FilingDate + ", " +
			patent.PublicationDate + ", " +
			patent.GrantDate + ", " +
			patent.Inventor + ", " +
			patent.Assignee + ", " +
			patent.Pdf
}

func PatentHeader() string {
	return `S/N, ApplicationNumber, Title, AssigneeOriginal, Language, PriorityDate, FilingDate, PublicationDate, GrantDate, Inventor, Assignee, Pdf`
}