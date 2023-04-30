// /*
// Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
// */

// package dispatcher

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"os"
// 	"strconv"
// 	"strings"
// 	"text/template"
// 	"time"

// 	sthingsBase "codehub.sva.de/Lab/stuttgart-things/dev/sthingsBase"
// 	sthingsK8s "codehub.sva.de/Lab/stuttgart-things/dev/sthingsK8s"
// 	redis "github.com/go-redis/redis/v7"
// 	redisqueue "github.com/robinjoseph08/redisqueue/v2"
// )

// var (
// 	redisAddress   = os.Getenv("REDIS_SERVER")
// 	redisPort      = os.Getenv("REDIS_PORT")
// 	redisPassword  = os.Getenv("REDIS_PASSWORD")
// 	yachtNamespace = os.Getenv("YACHT_NAMESPACE")
// 	log            = sthingsBase.StdOutFileLogger(logfilePath, "2006-01-02 15:04:05", 50, 3, 28)
// 	logfilePath    = "yad.log"
// 	redisClient    = redis.NewClient(&redis.Options{
// 		Addr:     redisAddress + ":" + redisPort,
// 		Password: redisPassword,
// 		DB:       0,
// 	})
// )

// type YawRevisionRunJob struct {
// 	Name                     string
// 	NamePrefix               string
// 	NameSuffix               string
// 	Image                    string
// 	TektonNamespace          string
// 	StatusToElasticsearch    string
// 	ElasticsearchUrl         string
// 	ElasticsearchStatusIndex string
// 	PrRanges                 string
// 	RevisionRunID            string
// }

// const jobTemplate = `
// apiVersion: batch/v1
// kind: Job
// metadata:
//   name: {{ .NamePrefix }}-{{ .Name }}-{{ .NameSuffix }}
//   namespace: yacht-tekton
//   labels:
//     jobgroup: yacht-application-worker
// spec:
//   template:
//     metadata:
//       name: yaw-3c5ac44c6fec00989c7e27b36630a82cdfd26e3b
//       namespace: yacht-tekton
//       labels:
//         jobgroup: yacht-application-worker
//     spec:
//       serviceAccountName: yacht-application-worker
//       containers:
//       - image: {{ .Image }}
//         envFrom:
//         - secretRef:
//             name: redis-connection
//         - configMapRef:
//             name: yaw-configuration
//         env:
//         - name: PR_RANGES
//           value: "{{ .PrRanges }}"
//         - name: REVISION_RUN_ID
//           value: {{ .RevisionRunID }}
//         - name: TEKTON_NAMESPACE
//           value: {{ .TektonNamespace }}
//         - name: STATUS_TO_ELASTICSEARCH
//           value: "{{ .StatusToElasticsearch }}"
//         - name: ELASTICSEARCH_URL
//           value: {{ .ElasticsearchUrl }}
//         - name: ELASTICSEARCH_STATUS_INDEX
//           value: {{ .ElasticsearchStatusIndex }}
//         resources:
//           limits:
//             cpu: 1000m
//             memory: 1000Mi
//           requests:
//             cpu: 500m
//             memory: 500Mi
//         name: yacht-application-worker
//       restartPolicy: Never
// `

// func CreateApplicationWorkerJobs(msg *redisqueue.Message) error {

// 	clusterConfig, clusterConnection := sthingsK8s.GetKubeConfig(os.Getenv("KUBECONFIG"))

// 	for _, element := range msg.Values {

// 		// Get revisionRun from redis q
// 		revisionRun := make(map[int][]string)
// 		pipelineRuns := fmt.Sprintf("%v", element)
// 		json.Unmarshal([]byte(pipelineRuns), &revisionRun)

// 		// create redis hash for revisionRun
// 		revisionRunMeta := make(map[string]interface{})
// 		yachtCommitID, _ := sthingsBase.GetRegexSubMatch(revisionRun[0][0], `yacht/commit: "(.*?)"`)
// 		fmt.Println("yachtCommitID:", yachtCommitID)

// 		//revisionRunMeta["ID"] = yachtCommitID
// 		log.Info("Processing revisionRun ", yachtCommitID, " right now")
// 		RedisAddHash(redisClient, "1", yachtCommitID, revisionRunMeta)
// 		log.Info("Added revisionRun to redis (hash: ", yachtCommitID, ")")

// 		log.Info("Connected " + clusterConnection + " the cluster")

// 		var prRanges []string

// 		for i := 0; i < (len(revisionRun)); i++ {

// 			prRanges = append(prRanges, sthingsBase.ConvertIntegerToString(len(revisionRun[i])))

// 			for j, pr := range revisionRun[i] {

// 				resourceName, _ := sthingsBase.GetRegexSubMatch(pr, `name: "(.*?)"`)
// 				fmt.Println("STAGE", i)
// 				fmt.Println(resourceName)
// 				RedisUpdateHashField(redisClient, yachtCommitID, strconv.Itoa(i)+":"+strconv.Itoa(j), pr)

// 			}

// 		}

// 		fmt.Println("RANGES", strings.Join(prRanges, ";"))
// 		SetRedisKeyValue(redisClient, "RANGE-"+yachtCommitID, strings.Join(prRanges, ";"))

// 		dt := time.Now()

// 		job := YawRevisionRunJob{
// 			Name:                     yachtCommitID,
// 			PrRanges:                 strings.Join(prRanges, ";"),
// 			NamePrefix:               "yaw",
// 			NameSuffix:               dt.Format("020405"),
// 			Image:                    os.Getenv("YAW_IMAGE"),
// 			TektonNamespace:          os.Getenv("TEKTON_NAMESPACE"),
// 			StatusToElasticsearch:    os.Getenv("STATUS_TO_ELASTICSEARCH"),
// 			ElasticsearchUrl:         os.Getenv("ELASTICSEARCH_URL"),
// 			ElasticsearchStatusIndex: os.Getenv("ELASTICSEARCH_STATUS_INDEX"),
// 			RevisionRunID:            yachtCommitID,
// 		}

// 		tmpl, err := template.New("pipelinerun").Parse(jobTemplate)
// 		if err != nil {
// 			panic(err)
// 		}

// 		var buf bytes.Buffer

// 		err = tmpl.Execute(&buf, job)

// 		if err != nil {
// 			log.Fatalf("execution: %s", err)
// 		}

// 		sthingsK8s.CreateDynamicResourcesFromTemplate(clusterConfig, []byte(buf.String()), yachtNamespace)

// 	}

// 	return nil
// }
