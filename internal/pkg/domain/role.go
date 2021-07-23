package domain

type Role struct {
	ID        int    `json:"id"`
	Num       *int   `json:"num,omitempty"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type AnimeStaffRole struct {
	Num     string `json:"num"`
	Role    string `json:"role"`
	Name    string `json:"name"`
	EngName string `json:"eng_name"`
}

type StaffRoleWithAnime struct {
	AnimeStaffRole
	TAnime
}

type RoleInteractor interface {
	StaffRoleByAnime(int) ([]AnimeStaffRole, error)
	// FilterByStaff(int) ([]StaffRoleWithAnime, error)
}
