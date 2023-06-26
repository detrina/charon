// Copyright © 2022-2023 Obol Labs Inc. Licensed under the terms of a Business Source License 1.1

package manifest_test

import (
	"encoding/json"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/obolnetwork/charon/cluster"
	"github.com/obolnetwork/charon/cluster/manifest"
	manifestpb "github.com/obolnetwork/charon/cluster/manifestpb/v1"
	"github.com/obolnetwork/charon/testutil"
)

//go:generate go test . -update

func TestZeroCluster(t *testing.T) {
	_, err := manifest.TypeLegacyLock.Transform(&manifestpb.Cluster{Name: "foo"}, &manifestpb.SignedMutation{})
	require.ErrorContains(t, err, "legacy lock not first mutation")
}

func TestLegacyLock(t *testing.T) {
	lockJON, err := os.ReadFile("testdata/lock.json")
	require.NoError(t, err)

	var lock cluster.Lock
	testutil.RequireNoError(t, json.Unmarshal(lockJON, &lock))

	signed, err := manifest.NewLegacyLock(lock)
	require.NoError(t, err)

	t.Run("proto", func(t *testing.T) {
		testutil.RequireGoldenProto(t, signed)
	})

	t.Run("cluster", func(t *testing.T) {
		cluster, err := manifest.Materialise(&manifestpb.SignedMutationList{Mutations: []*manifestpb.SignedMutation{signed}})
		require.NoError(t, err)
		require.Equal(t, lock.LockHash, cluster.Hash)
		testutil.RequireGoldenProto(t, cluster)
	})

	b, err := proto.Marshal(signed)
	require.NoError(t, err)

	signed2 := new(manifestpb.SignedMutation)
	testutil.RequireNoError(t, proto.Unmarshal(b, signed2))

	t.Run("proto again", func(t *testing.T) {
		testutil.RequireGoldenProto(t, signed2, testutil.WithFilename("TestLegacyLock_proto.golden"))
	})

	t.Run("cluster loaded from lock", func(t *testing.T) {
		cluster, err := manifest.Load("testdata/lock.json", nil)
		require.NoError(t, err)
		testutil.RequireGoldenProto(t, cluster, testutil.WithFilename("TestLegacyLock_cluster.golden"))
	})

	t.Run("cluster loaded from manifest", func(t *testing.T) {
		b, err := proto.Marshal(&manifestpb.SignedMutationList{Mutations: []*manifestpb.SignedMutation{signed}})
		require.NoError(t, err)
		file := path.Join(t.TempDir(), "manifest.pb")
		require.NoError(t, os.WriteFile(file, b, 0o644))

		cluster, err := manifest.Load(file, nil)
		require.NoError(t, err)
		testutil.RequireGoldenProto(t, cluster, testutil.WithFilename("TestLegacyLock_cluster.golden"))
	})
}