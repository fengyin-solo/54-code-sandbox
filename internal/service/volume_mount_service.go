package service

import (
	"sort"
	"time"

	"sandbox/internal/model"
	"sandbox/pkg/idgen"
)

func (s *Service) CreateVolumeMount(input model.VolumeMount) (*model.VolumeMount, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetSandbox(input.SandboxID); err != nil {
		return nil, model.NewValidationError("sandbox_id", "所属沙箱不存在")
	}
	now := time.Now()
	v := &model.VolumeMount{
		ID:            idgen.Hex(),
		SandboxID:     input.SandboxID,
		HostPath:      input.HostPath,
		ContainerPath: input.ContainerPath,
		ReadOnly:      input.ReadOnly,
		Status:        input.Status,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if err := s.store.CreateVolumeMount(v); err != nil {
		return nil, err
	}
	s.log.Infof("创建卷挂载: %s", v.ID)
	return v, nil
}

func (s *Service) GetVolumeMount(id string) (*model.VolumeMount, error) {
	return s.store.GetVolumeMount(id)
}

func (s *Service) ListVolumeMounts(filter model.VolumeMountFilter, page, size int) ([]*model.VolumeMount, int, error) {
	all := s.store.ListVolumeMounts()
	matched := make([]*model.VolumeMount, 0, len(all))
	for _, v := range all {
		if filter.Match(v) {
			matched = append(matched, v)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.VolumeMount{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateVolumeMount(id string, input model.VolumeMount) (*model.VolumeMount, error) {
	v, err := s.store.GetVolumeMount(id)
	if err != nil {
		return nil, err
	}
	if input.HostPath != "" {
		v.HostPath = input.HostPath
	}
	if input.ContainerPath != "" {
		v.ContainerPath = input.ContainerPath
	}
	v.ReadOnly = input.ReadOnly
	if input.Status != "" {
		v.Status = input.Status
	}
	v.UpdatedAt = time.Now()
	if err := v.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateVolumeMount(v); err != nil {
		return nil, err
	}
	return v, nil
}

func (s *Service) DeleteVolumeMount(id string) error {
	return s.store.DeleteVolumeMount(id)
}
