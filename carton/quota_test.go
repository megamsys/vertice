package carton

/*
import (
	"gopkg.in/check.v1"
	"strconv"
)

func (s *S) TestGetQuota(c *check.C) {
	q := new(Quota)
	q.AccountId = "mvijaykanth@megam.io"
	q.Id = "QUO8671194369186000000"
	s.Credentials.Email = q.AccountId
  _, err := q.get(s.Credentials)
  c.Assert(err, check.IsNil)
}

func (s *S) TestUpdateQuota(c *check.C) {
	q := new(Quota)
	q.AccountId = "mvijaykanth@megam.io"
	q.Id = "QUO8671194369186000000"
	s.Credentials.Email = q.AccountId
	res, err := q.get(s.Credentials)
	c.Assert(err, check.IsNil)
	val := res.AllowedSnaps()
	count, _ := strconv.Atoi(res.AllowedSnaps())
  mm := make(map[string][]string, 1)
  mm["snapshots"] = []string{strconv.Itoa(count-1)}
  res.Allowed.NukeAndSet(mm)
	err = res.update(s.Credentials)
	c.Assert(nil, check.NotNil)
}
*/
