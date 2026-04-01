package interfaces

type RoomUsecase struct {
	jwt JWTService
	repo RoomRepository
}

func NewRoomUsecase(jwt JWTService, repo RoomRepository) *RoomUsecase {
	return &RoomUsecase{jwt: jwt, repo: repo}
}

func (r *RoomUsecase) ListRooms() OutputDTO {
	list, err := r.repo.GetAllRooms()
	if (err != nil) {
		return OutputDTO{Status: 500, Data: map[string]interface{}{"error": err.Error()}}
	}
	return OutputDTO{Status: 200, Data: map[string]interface{}{"list": *list}}
}

func (r *RoomUsecase) CreateRoom(input CreateDTO, token string) OutputDTO {
	user, err := r.jwt.ValidateToken(token)
	if (err != nil) {
		return OutputDTO{Status: 500, Data: map[string]interface{}{"error": err.Error()}}
	}
	if (user.Role != "admin") {
		return OutputDTO{Status: 403, Data: map[string]interface{}{"error": "for this action need admin role"}}
	}
	room, err := r.repo.CreateRoom(input.Name, input.Description, input.Capacity)
	if (err != nil) {
		return OutputDTO{Status: 400, Data: map[string]interface{}{"error": err.Error()}}
	}
	return OutputDTO{Status: 201, Data: map[string]interface{}{"room": *room}}
}

func (r *RoomUsecase) GetRoom(input GetDTO) OutputDTO {
	room, err := r.repo.GetRoom(input.Name)
	if (err != nil) {
		return OutputDTO{Status: 400, Data: map[string]interface{}{"error": err.Error()}}
	}
	return OutputDTO{Status: 200, Data: map[string]interface{}{"list": *room}}
}