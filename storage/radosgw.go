package storage

import (
  "github.com/megamsys/go-radosgw/api"
  //"strconv"
)

type RadosGW struct {
  UserId string `json:"uid"`
  Api *radosAPI.API
  TotalSizeMB float64 `json:"total_size"`
}




func NewRgW(host,access, secret string) *RadosGW {
  return &RadosGW{Api: radosAPI.New(host, access, secret)}
}

// func (r *RadosGW) getUsers(user string) (*User, error) {
// 	user, err := r.Api.GetUser(u)
// 	if err != nil {
//     return nil, err
// 	}
//   return r.userParse(user)
// }

// func (r *RadosGW) userParse(*radosAPI.User) (*User, error) {
//   // have to parse rados user to User
//   return new(User), nil
// }

// returns user's storage size of all buckets
func (r *RadosGW) getUserBuckets(name, user string) (radosAPI.Buckets, error) {
   	bkt := radosAPI.BucketConfig{Bucket: name, UID: user, Stats: true}
   	buckets, err := r.Api.GetBuckets(bkt)
	if err != nil {
		return nil, err
	}
  return buckets, nil
}

func (r *RadosGW) totalSize(buckets radosAPI.Buckets) error {
   var size float64 = 0.0
    for _, b := range buckets {
      size = size + float64(float64(b.Stats.Usage.RgwMain.SizeKbActual)/1024.0)
    }
  r.TotalSizeMB = size
  return nil
}

func (r *RadosGW) GetUserStorageSize(user string) error {
  buckets, err  := r.getUserBuckets("",user)
  if err != nil {
    return  err
  }
 return r.totalSize(buckets)
}

func (r *RadosGW) GetUserBucketSize(name, user string) error {
  buckets, err  := r.getUserBuckets(name, user)
  if err != nil {
    return  err
  }
 return r.totalSize(buckets)
}
