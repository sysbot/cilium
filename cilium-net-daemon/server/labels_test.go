package server

import (
	"errors"
	"strings"

	"github.com/noironetworks/cilium-net/common"
	"github.com/noironetworks/cilium-net/common/types"

	. "github.com/noironetworks/cilium-net/Godeps/_workspace/src/gopkg.in/check.v1"
)

var (
	lbls = createLbls()

	wantSecCtxLbls = types.SecCtxLabels{
		ID:       123,
		RefCount: 1,
		Labels:   lbls,
	}
)

func createLbls() types.Labels {
	lbls := []types.Label{
		types.NewLabel("foo", "bar", common.CiliumLabelSource),
		types.NewLabel("foo2", "=bar2", common.CiliumLabelSource),
		types.NewLabel("key", "", common.CiliumLabelSource),
		types.NewLabel("foo==", "==", common.CiliumLabelSource),
		types.NewLabel(`foo\\=`, `\=`, common.CiliumLabelSource),
		types.NewLabel(`//=/`, "", common.CiliumLabelSource),
		types.NewLabel(`%`, `%ed`, common.CiliumLabelSource),
	}
	return map[string]*types.Label{
		"foo":    &lbls[0],
		"foo2":   &lbls[1],
		"key":    &lbls[2],
		"foo==":  &lbls[3],
		`foo\\=`: &lbls[4],
		`//=/`:   &lbls[5],
		`%`:      &lbls[6],
	}
}

func (s *DaemonSuite) TestGetLabelsIDOK(c *C) {
	s.d.OnPutLabels = func(lblsReceived types.Labels) (*types.SecCtxLabels, bool, error) {
		c.Assert(lblsReceived, DeepEquals, lbls)
		return &wantSecCtxLbls, true, nil
	}

	secCtxLabl, _, err := s.c.PutLabels(lbls)
	c.Assert(err, Equals, nil)
	c.Assert(*secCtxLabl, DeepEquals, wantSecCtxLbls)
}

func (s *DaemonSuite) TestGetLabelsIDFail(c *C) {
	s.d.OnPutLabels = func(lblsReceived types.Labels) (*types.SecCtxLabels, bool, error) {
		c.Assert(lblsReceived, DeepEquals, lbls)
		return nil, false, errors.New("Reached maximum valid IDs")
	}

	_, _, err := s.c.PutLabels(lbls)
	c.Assert(strings.Contains(err.Error(), "Reached maximum valid IDs"), Equals, true)
}

func (s *DaemonSuite) TestGetLabelsOK(c *C) {
	s.d.OnGetLabels = func(id int) (*types.SecCtxLabels, error) {
		c.Assert(id, Equals, 123)
		return &wantSecCtxLbls, nil
	}

	lblsReceived, err := s.c.GetLabels(123)
	c.Assert(err, Equals, nil)
	c.Assert(*lblsReceived, DeepEquals, wantSecCtxLbls)
}

func (s *DaemonSuite) TestGetLabelsFail(c *C) {
	s.d.OnGetLabels = func(id int) (*types.SecCtxLabels, error) {
		c.Assert(id, Equals, 123)
		return nil, errors.New("Unable to contact consul")
	}

	_, err := s.c.GetLabels(123)
	c.Assert(strings.Contains(err.Error(), "Unable to contact consul"), Equals, true)
}

func (s *DaemonSuite) TestDeleteLabelsByUUIDOK(c *C) {
	s.d.OnDeleteLabelsByUUID = func(id int) error {
		c.Assert(id, Equals, 123)
		return nil
	}

	err := s.c.DeleteLabelsByUUID(123)
	c.Assert(err, Equals, nil)
}

func (s *DaemonSuite) TestDeleteLabelsByUUIDFail(c *C) {
	s.d.OnDeleteLabelsByUUID = func(id int) error {
		c.Assert(id, Equals, 123)
		return errors.New("Unable to contact consul")
	}

	err := s.c.DeleteLabelsByUUID(123)
	c.Assert(strings.Contains(err.Error(), "Unable to contact consul"), Equals, true)
}

func (s *DaemonSuite) TestDeleteLabelsBySHA256OK(c *C) {
	s.d.OnDeleteLabelsBySHA256 = func(sha256sum string) error {
		c.Assert(sha256sum, Equals, "82078f981c61a5a71acbe92d38b2de3e3c5f7469450feab03d2739dfe6cbc049")
		return nil
	}

	err := s.c.DeleteLabelsBySHA256("82078f981c61a5a71acbe92d38b2de3e3c5f7469450feab03d2739dfe6cbc049")
	c.Assert(err, Equals, nil)
}

func (s *DaemonSuite) TestDeleteLabelsBySHA256Fail(c *C) {
	s.d.OnDeleteLabelsBySHA256 = func(sha256sum string) error {
		c.Assert(sha256sum, Equals, "82078f981c61a5a71acbe92d38b2de3e3c5f7469450feab03d2739dfe6cbc049")
		return errors.New("Unable to contact consul")
	}

	err := s.c.DeleteLabelsBySHA256("82078f981c61a5a71acbe92d38b2de3e3c5f7469450feab03d2739dfe6cbc049")
	c.Assert(strings.Contains(err.Error(), "Unable to contact consul"), Equals, true)
}

func (s *DaemonSuite) TestGetMaxOK(c *C) {
	s.d.OnGetMaxID = func() (int, error) {
		return 100, nil
	}

	maxID, err := s.c.GetMaxID()
	c.Assert(err, Equals, nil)
	c.Assert(maxID, Equals, 100)
}

func (s *DaemonSuite) TestGetMaxIDFail(c *C) {
	s.d.OnGetMaxID = func() (int, error) {
		return -1, errors.New("Unable to contact consul")
	}

	_, err := s.c.GetMaxID()
	c.Assert(strings.Contains(err.Error(), "Unable to contact consul"), Equals, true, Commentf("error %s", err))
}
