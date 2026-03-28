package helper

import (
	"sort"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const instanceLabelKey = "neptune.io/instance"

// BuildSchedulingAffinity builds node constraints from planned nodes and
// spreads same-instance pods across hosts when possible.
func BuildSchedulingAffinity(instanceName string, allowedNodes []string, strictSpread bool) *corev1.Affinity {
	affinity := &corev1.Affinity{}

	if nodes := uniqueNodeNames(allowedNodes); len(nodes) > 0 {
		affinity.NodeAffinity = &corev1.NodeAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
				NodeSelectorTerms: []corev1.NodeSelectorTerm{
					{
						MatchExpressions: []corev1.NodeSelectorRequirement{
							{
								Key:      corev1.LabelHostname,
								Operator: corev1.NodeSelectorOpIn,
								Values:   nodes,
							},
						},
					},
				},
			},
		}
	}

	if instanceName != "" {
		term := corev1.PodAffinityTerm{
			LabelSelector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					instanceLabelKey: instanceName,
				},
			},
			TopologyKey: corev1.LabelHostname,
		}

		if strictSpread {
			affinity.PodAntiAffinity = &corev1.PodAntiAffinity{
				RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{term},
			}
		} else {
			affinity.PodAntiAffinity = &corev1.PodAntiAffinity{
				PreferredDuringSchedulingIgnoredDuringExecution: []corev1.WeightedPodAffinityTerm{
					{
						Weight:          100,
						PodAffinityTerm: term,
					},
				},
			}
		}
	}

	if affinity.NodeAffinity == nil && affinity.PodAntiAffinity == nil {
		return nil
	}

	return affinity
}

func uniqueNodeNames(nodes []string) []string {
	if len(nodes) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(nodes))
	result := make([]string, 0, len(nodes))
	for _, node := range nodes {
		if node == "" {
			continue
		}
		if _, ok := seen[node]; ok {
			continue
		}
		seen[node] = struct{}{}
		result = append(result, node)
	}

	sort.Strings(result)
	return result
}
