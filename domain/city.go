package domain

import "context"

type City struct {
	Code      int     `db:"code"`
	Province  string  `db:"province"`
	City      string  `db:"city"`
	Latitude  float32 `db:"latitude"`
	Longitude float32 `db:"longitude"`
	Timezone  int     `db:"timezone"`
	IDmyQuran int     `db:"id_my_quran"`
}

type CityResponseDTO struct {
	CityCode  int     `json:"city_code"`
	City      string  `json:"city"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
	Timezone  int     `json:"timezone"`
}

type CityUsecase interface {
	GetByLatLong(ctx context.Context, latitude float32, longitude float32) (res CityResponseDTO, err error)
	GetByCityCode(ctx context.Context, cityCode int) (res CityResponseDTO, err error)
	// GetAll(c context.Context, offset, limit int64, sort, province string, desc string) ([]MissIndonesiaListResponseDTO, int64, error)
	// GetParent(ctx context.Context, userID string) (*ParentData, error)
	// GetByUserID(ctx context.Context, userID string) (MissIndonesiaResponseDTO, error)
	// Create(ctx context.Context, userID string, mi *MissIndonesiaRequestDTO) error
	// UploadPhoto(ctx context.Context, filephoto []byte, tipe string, userID string) error
	// UploadVideo(ctx context.Context, video []byte, typeVid, inputVid string) error
	// PaymentOrder(c context.Context, po *PaymentOrderRequestDTO, id string) (*PaymentOrderResponseDTO, string, error)
	// Update(ctx context.Context, userID string, mi *MissIndonesiaRequestUpdateDTO) error
	// Delete(ctx context.Context, userID string) error
	// GetUserByIDAdmin(ctx context.Context, userID string) (res MissIndonesiaListResponseDTO, err error)
	// InsertToXLSX(ctx context.Context, province string, sort, desc string) error
	// HookUpload(ctx context.Context, request *GRPChook.SendRequest) error
	// Validation(ctx context.Context, request *GRPChook.SendRequest) error
}

type CityRepository interface {
	Search(ctx context.Context, offset, limit int64, sort, province string, desc string) (res []City, err error)
	GetByLatLong(ctx context.Context, latitude float32, longitude float32) (res City, err error)
	GetByCityCode(ctx context.Context, cityCode int) (res City, err error)
	GetByCityName(ctx context.Context, cityName string) (res City, err error)
	Create(ctx context.Context, city *City) (err error)
	Update(ctx context.Context, city *City) (err error)
	Delete(ctx context.Context, CityCode string) (err error)
}
