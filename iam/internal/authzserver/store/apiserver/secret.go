package apiserver

import (
	"context"
	"github.com/AlekSi/pointer"
	"github.com/avast/retry-go/v4"
	pb "istomyang.github.com/like-iam/api/proto/v1"
	"istomyang.github.com/like-iam/component-base/errors"
	"istomyang.github.com/like-iam/iam/internal/authzserver/store"
	"istomyang.github.com/like-iam/log"
)

type secret struct {
	pb  pb.CacheClient
	ctx context.Context
}

func newSecret(ctx context.Context, pb pb.CacheClient) store.SecretStore {
	return &secret{pb: pb, ctx: ctx}
}

func (s *secret) List() (map[string]*pb.SecretInfo, error) {

	log.Info("loading list secrets.")

	req := pb.ListRequest{
		Offset: pointer.ToInt64(0),
		Limit:  pointer.ToInt64(-1), // cancel offset condition with -1
	}

	var secrets *pb.ListSecretsReply
	var err error

	err = retry.Do(func() error {
		secrets, err = s.pb.ListSecrets(s.ctx, &req)
		return err
	}, retry.Attempts(3))
	if err != nil {
		return nil, errors.Wrap(err, "list secrets from apiserver failed after 3 times.")
	}

	log.Infof("secrets loaded count: &d", secrets.Count)

	r := make(map[string]*pb.SecretInfo, secrets.Count)

	for _, item := range secrets.Items {
		log.Infof("get secrets: %s:%s", item.Username, item.SecretId)
		r[item.SecretId] = item
	}

	return r, nil
}
