/*
Copyright 2022 The KubeService-Stack Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package framework

import (
	"context"

	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

var (
	DefaultCSIRule = rbacv1.PolicyRule{
		APIGroups: []string{"csi.aliyun.com"},
		Resources: []string{
			"nodelocalstorages",
			"nodelocalstorages/status",
			"nodelocalstorageinitconfigs",
		},
		Verbs: []string{
			"create",
			"get",
			"list",
			"watch",
			"update",
			"delete",
			"patch",
		},
	}

	DefaultEventsRule = rbacv1.PolicyRule{
		APIGroups: []string{""},
		Resources: []string{"events"},
		Verbs:     []string{"create", "update", "patch"},
	}

	DefaultCoreRule = rbacv1.PolicyRule{
		APIGroups: []string{""},
		Resources: []string{
			"nodes",
			"pods",
			"pods/binding",
			"pods/status",
			"bindings",
			"persistentvolumeclaims",
			"persistentvolumeclaims/status",
			"persistentvolumes",
			"persistentvolumes/status",
			"namespaces",
			"secrets",
		},
		Verbs: []string{"create", "get", "list", "watch", "update", "delete", "patch"},
	}

	DefaultStorageRule = rbacv1.PolicyRule{
		APIGroups: []string{"storage.k8s.io"},
		Resources: []string{
			"storageclasses",
			"csinodes",
			"volumeattachments",
		},
		Verbs: []string{"get", "list", "watch"},
	}

	DefaultSnapshotRule = rbacv1.PolicyRule{
		APIGroups: []string{"snapshot.storage.k8s.io"},
		Resources: []string{
			"volumesnapshotclasses",
			"volumesnapshots",
			"volumesnapshots/status",
			"volumesnapshotcontents",
			"volumesnapshotcontents/status",
		},
		Verbs: []string{"create", "get", "list", "watch", "update", "delete", "patch"},
	}

	DefaultCoordinationRule = rbacv1.PolicyRule{
		APIGroups: []string{"coordination.k8s.io"},
		Resources: []string{
			"leases",
		},
		Verbs: []string{"create", "get", "list", "watch", "update", "delete", "patch"},
	}
)

func (f *Framework) CreateOrUpdateClusterRole(ctx context.Context, source string) (*rbacv1.ClusterRole, error) {
	clusterRole, err := parseClusterRoleYaml(source)
	if err != nil {
		return nil, err
	}

	_, err = f.KubeClient.RbacV1().ClusterRoles().Get(ctx, clusterRole.Name, metav1.GetOptions{})
	if err != nil && !apierrors.IsNotFound(err) {
		return nil, err
	}

	if apierrors.IsNotFound(err) {
		// ClusterRole doesn't exists -> Create
		clusterRole, err = f.KubeClient.RbacV1().ClusterRoles().Create(ctx, clusterRole, metav1.CreateOptions{})
		if err != nil {
			return nil, err
		}
	} else {
		// ClusterRole already exists -> Update
		clusterRole, err = f.KubeClient.RbacV1().ClusterRoles().Update(ctx, clusterRole, metav1.UpdateOptions{})
		if err != nil {
			return nil, err
		}
	}

	return clusterRole, nil
}

func (f *Framework) DeleteClusterRole(ctx context.Context, source string) error {
	clusterRole, err := parseClusterRoleYaml(source)
	if err != nil {
		return err
	}

	return f.KubeClient.RbacV1().ClusterRoles().Delete(ctx, clusterRole.Name, metav1.DeleteOptions{})
}

func (f *Framework) UpdateClusterRole(ctx context.Context, clusterRole *rbacv1.ClusterRole) error {
	_, err := f.KubeClient.RbacV1().ClusterRoles().Update(ctx, clusterRole, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func parseClusterRoleYaml(source string) (*rbacv1.ClusterRole, error) {
	manifest, err := SourceToIOReader(source)
	if err != nil {
		return nil, err
	}

	clusterRole := rbacv1.ClusterRole{}
	if err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&clusterRole); err != nil {
		return nil, err
	}

	return &clusterRole, nil
}
