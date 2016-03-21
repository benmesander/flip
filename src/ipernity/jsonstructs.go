package ipernity

type Userget struct {
	User struct {
		User_id   string
		Is_pro    string
		Is_online string
		Is_closed string
		Link      string
		Username  string
		Realname  string
		Sex       string
		Count     struct {
			Network string
			Docs    string
			Posts   string
		}
		Location struct {
			Country string
			Town    string
		}
		Dates struct {
			Last_online  string
			Member_since string
		}
	}
	Api struct {
		Status  string `json:"status"`
		At      string `json:"at"`
		Code    string `json:"code"`    // only if status != "ok"
		Message string `json:"message"` // only if status != ok
	}
}

type Docgetlist struct {
	Docs struct {
		Total    string
		Per_page string
		Page     string
		Pages    string
		Count    string

		Doc []struct {
			Doc_id   string
			Media    string
			Title    string
			Original struct {
				Url   string
				Bytes string
			}
		}
	}
	Api struct {
		Status  string `json:"status"`
		At      string `json:"at"`
		Code    string `json:"code"`    // only if status != "ok"
		Message string `json:"message"` // only if status != ok
	}
}

type Docget struct {
	Doc struct {
		Can struct {
			Comment   string
			Download  string
			Extra     string
			Fave      string
			Print     string
			Share     string
			Tag       string
			Tagme     string
			Translate string
			Zoom      string
		}
		Count struct {
			Albums   string
			Comments string
			Faves    string
			Groups   string
			Notes    string
			Tags     string
			Visits   string
		}
		Dates struct {
			Created         string
			Last_comment_at string
			Last_update_at  string
			Posted_at       string
		}
		Description string
		Doc_id      string
		Icon        string
		License     string
		Media       string
		Original    struct {
			Bytes    string
			Ext      string
			Filename string
			H        string
			Url      string
			Width    string
		}
		Owner struct {
			Alias    string
			Icon     string
			Is_pro   string
			User_id  string
			Username string
		}
		Permissions struct {
			Comment string
			Tag     string
			Tagme   string
		}
		Rotation string
		Share    struct {
			Url string
		}
		Thumbs struct {
			Thumb []struct {
				Ext    string
				Farm   string
				H      string
				Icon   string
				Label  string
				Path   string
				Secret string
				Url    string
				W      string
			}
		}
		Title      string
		Visibility struct {
			Isfamily string
			Isfriend string
			Ispublic string
			Share    string
		}
		You struct {
			Isfave        string
			Last_visit_at string
			Visits        string
		}
	}

	Api struct {
		Status  string `json:"status"`
		At      string `json:"at"`
		Code    string `json:"code"`    // only if status != "ok"
		Message string `json:"message"` // only if status != ok
	}
}

type Docgetcontainers struct {
	Albums struct {
		Album []struct {
			Album_id string
			Cover    string
			Docs     string
			Link     string
			Title    string
		}
		Total string
	}

	Groups struct {
		Group []struct {
			Cover    string
			Docs     string
			Group_id string
			Link     string
			Title    string
		}
		Total string
	}

	Api struct {
		Status  string `json:"status"`
		At      string `json:"at"`
		Code    string `json:"code"`    // only if status != "ok"
		Message string `json:"message"` // only if status != ok
	}
}
