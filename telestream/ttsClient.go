package telestream

import (
	"context"
	"fmt"
	"tcs-cli/cli"
	"github.com/Telestream/telestream-cloud-go-sdk/tts"
)

// client holds output writer, tts configuration, tts client and context
type TtsClient struct {
	config *tts.Configuration
	client *tts.APIClient
	ctx    context.Context
	output ServiceOutput
}

// Creates new telestream tts client with given output writer and X API key
func NewTtsClient(xApiKey string, addHearderKey string, addHeaderVal string, output ServiceOutput) *TtsClient {

	client := new(TtsClient)
	client.config = tts.NewConfiguration()

	if addHearderKey != "" && addHeaderVal != "" {

		client.config.AddDefaultHeader(addHearderKey, addHeaderVal)
	}

	client.client = tts.NewAPIClient(client.config)
	client.ctx = context.WithValue(context.Background(), tts.ContextAPIKey, tts.APIKey{Key: xApiKey})
	client.output = output

	return client
}

// List all projets to output
func (client *TtsClient) ListProjects() {

	projectsCollection, _, err := client.client.TtsApi.Projects(client.ctx)

	colNames := []interface{}{"NAME", "ID", "CREATED_AT", "STATUS", "DESCRIPTION"}
	rows := [][]interface{}{}

	if nil == err {

		for _, project := range projectsCollection.Projects {

			rows = append(rows, []interface{}{project.Name, project.Id, project.CreatedAt,
				project.Status, project.Description})
		}

		client.output.printTable(colNames, rows)

	} else {

		client.output.printError("ListProjects", err)
	}
}

// Print project description on output
func (client *TtsClient) DescribeProject(argsMap cli.FlagMap) {

	projectDesc, _, err := client.client.TtsApi.Project(client.ctx, *argsMap["project_id"].Value)

	if nil == err {

		client.output.printStructContent(&projectDesc)

	} else {

		client.output.printError("DescribeProject", err)
	}
}

// Get describe factory input attributes
func (client *TtsClient) GetDescribeProjectProperties() map[string]bool {

	flagMap := map[string]bool{"project_id": true}

	return flagMap
}

// Create new project
func (client *TtsClient) CreateProject(argsMap cli.FlagMap) {

	newProject := tts.Project{}
	propertiesToStruct(&newProject, argsMap)

	projectDesc, _, err := client.client.TtsApi.CreateProject(client.ctx, newProject)

	if nil == err {

		client.output.printStructContent(&projectDesc)

	} else {

		client.output.printError("CreateProject", err)
	}
}

// Get create project input attributes
func (client *TtsClient) GetCreateProjectProperties() map[string]bool {

	flagMap := map[string]bool{}
	project := tts.Project{}

	jsonFields := structToProperties(&project)

	for _, field := range jsonFields {

		flagMap[field] = false
	}

	flagMap["name"] = true
	flagMap["description"] = true
	flagMap["language"] = true

	return flagMap
}

// Delete project
func (client *TtsClient) DeleteProject(argsMap cli.FlagMap) {

	_, err := client.client.TtsApi.DeleteProject(client.ctx, *argsMap["project_id"].Value)

	if nil == err {

		client.output.printInfo("Project: " + *argsMap["project_id"].Value + " removed.")

	} else {

		client.output.printError("DescribeProject", err)
	}
}

// Get delete project input attributes
func (client *TtsClient) GetDeleteProjectProperties() map[string]bool {

	flagMap := map[string]bool{"project_id": true}

	return flagMap
}

// Update project and print its new body
func (client *TtsClient) UpdateProject(argsMap cli.FlagMap) {

	newProject := tts.Project{}
	propertiesToStruct(&newProject, argsMap)

	projectDesc, _, err := client.client.TtsApi.UpdateProject(client.ctx, newProject.Id, newProject)

	if nil == err {

		client.output.printStructContent(&projectDesc)

	} else {

		client.output.printError("UpdateProject", err)
	}
}

// Get update project input attributes
func (client *TtsClient) GetUpdateProjectProperties() map[string]bool {

	flagMap := map[string]bool{}
	project := tts.Project{}

	jsonFields := structToProperties(&project)

	for _, field := range jsonFields {

		flagMap[field] = false
	}

	flagMap["id"] = true

	return flagMap
}

// List all projets to output
func (client *TtsClient) ListJobs(argsMap cli.FlagMap) {

	opts, pageErr := getPageOpt(&argsMap)
	if pageErr != nil {

		client.output.printError("ListJobs", pageErr)
		return
	}

	jobsCollection, _, err := client.client.TtsApi.Jobs(client.ctx, *argsMap["project_id"].Value,
		opts)

	// print all projects in table
	colNames := []interface{}{"JOB_ID", "CREATED_AT", "STATUS", "STREAM_NAME", "DURATION", "CONFIDENCE"}
	rows := [][]interface{}{}

	if nil == err {

		for _, job := range jobsCollection.Jobs {

			rows = append(rows, []interface{}{job.Id, job.CreatedAt, job.Status, job.Name, fmt.Sprint(job.Duration), fmt.Sprint(job.Confidence)})
		}

		client.output.printTable(colNames, rows)
	} else {

		client.output.printError("ListJobs", err)
	}
}

// Get list jobs input attributes
func (client *TtsClient) GetListJobsProperties() map[string]bool {

	flagMap := map[string]bool{"project_id": true}
	addPageOpt(flagMap)

	return flagMap
}

// Create new job
func (client *TtsClient) CreateJob(argsMap cli.FlagMap) {

	newJob := tts.Job{}
	propertiesToStruct(&newJob, argsMap)

	jobDesc, _, err := client.client.TtsApi.CreateJob(client.ctx, newJob.ProjectId, newJob)

	if nil == err {

		client.output.printStructContent(&jobDesc)

	} else {

		client.output.printError("CreateJob", err)
	}
}

// Get create job input attributes
func (client *TtsClient) GetCreateJobProperties() map[string]bool {

	flagMap := map[string]bool{}
	job := tts.Job{}

	jsonFields := structToProperties(&job)

	for _, field := range jsonFields {

		flagMap[field] = false
	}

	flagMap["project_id"] = true
	flagMap["source_url"] = true

	return flagMap
}

// Print job description on output
func (client *TtsClient) DescribeJob(argsMap cli.FlagMap) {

	jobDesc, _, err := client.client.TtsApi.Job(client.ctx, *argsMap["project_id"].Value,
		*argsMap["job_id"].Value)

	if nil == err {

		client.output.printStructContent(&jobDesc)

	} else {

		client.output.printError("DescribeJob", err)
	}
}

// Get describe job input attributes
func (client *TtsClient) GetDescribeJobProperties() map[string]bool {

	flagMap := map[string]bool{"project_id": true, "job_id": true}

	return flagMap
}

// Delete job
func (client *TtsClient) DeleteJob(argsMap cli.FlagMap) {

	_, err := client.client.TtsApi.DeleteJob(client.ctx, *argsMap["project_id"].Value,
		*argsMap["job_id"].Value)

	if nil == err {

		client.output.printInfo("Deleted job: " + *argsMap["job_id"].Value + " in project: " +
			*argsMap["project_id"].Value)

	} else {

		client.output.printError("DeleteJob", err)
	}
}

// Get delete job input attributes
func (client *TtsClient) GetDeleteJobProperties() map[string]bool {

	flagMap := map[string]bool{"project_id": true, "job_id": true}

	return flagMap
}

// Print job result on output
func (client *TtsClient) JobResult(argsMap cli.FlagMap) {

	jobResult, _, err := client.client.TtsApi.JobResult(client.ctx, *argsMap["project_id"].Value,
		*argsMap["job_id"].Value)

	if nil == err {

		client.output.printStructContent(&jobResult)

	} else {

		client.output.printError("JobResult", err)
	}
}

// Get job result input attributes
func (client *TtsClient) GetJobResultProperties() map[string]bool {

	flagMap := map[string]bool{"project_id": true, "job_id": true}

	return flagMap
}

// Print job output on output
func (client *TtsClient) JobOutputs(argsMap cli.FlagMap) {

	jobOutputs, _, err := client.client.TtsApi.JobOutputs(client.ctx, *argsMap["project_id"].Value,
		*argsMap["job_id"].Value)

	if nil == err {

		for _, output := range jobOutputs {

			client.output.printStructContent(&output)
		}

	} else {

		client.output.printError("JobResult", err)
	}
}

// Get job result input attributes
func (client *TtsClient) GetJobOutputsProperties() map[string]bool {

	flagMap := map[string]bool{"project_id": true, "job_id": true}

	return flagMap
}

// List all corpora
func (client *TtsClient) ListCorpora(argsMap cli.FlagMap) {

	corporaCollection, _, err := client.client.TtsApi.Corpora(client.ctx, *argsMap["project_id"].Value)

	// print all corpora in table
	colNames := []interface{}{"NAME", "STATUS"}
	rows := [][]interface{}{}

	if nil == err {

		for _, corpus := range corporaCollection.Corpora {
			rows = append(rows, []interface{}{corpus.Name, corpus.Status})
		}

		client.output.printTable(colNames, rows)

	} else {

		client.output.printError("ListCorpora", err)
	}
}

// Get list corpora input attributes
func (client *TtsClient) GetListCorporaProperties() map[string]bool {

	flagMap := map[string]bool{"project_id": true}

	return flagMap
}

// Print corpus description on output
func (client *TtsClient) DescribeCorpus(argsMap cli.FlagMap) {

	corpusDesc, _, err := client.client.TtsApi.Corpus(client.ctx, *argsMap["project_id"].Value,
		*argsMap["corpus_name"].Value)

	if nil == err {

		client.output.printStructContent(&corpusDesc)

	} else {

		client.output.printError("DescribeCorpus", err)
	}
}

// Get describe corpus input attributes
func (client *TtsClient) GetDescribeCorpusProperties() map[string]bool {

	flagMap := map[string]bool{"project_id": true, "corpus_name": true}

	return flagMap
}

// Create new corpus
func (client *TtsClient) CreateCorpus(argsMap cli.FlagMap) {

	_, err := client.client.TtsApi.CreateCorpus(client.ctx, *argsMap["project_id"].Value,
		*argsMap["corpus_name"].Value, *argsMap["corpus_body"].Value)

	if nil == err {

		client.output.printInfo("Corpus created: " + *argsMap["corpus_name"].Value + " in project: " +
			*argsMap["project_id"].Value)

	} else {

		client.output.printError("CreateCorpus ", err)
	}
}

// Get create corpus input attributes
func (client *TtsClient) GetCreateCorpusProperties() map[string]bool {

	flagMap := map[string]bool{"project_id": true, "corpus_name": true, "corpus_body": true}

	return flagMap
}

// Delete corpus
func (client *TtsClient) DeleteCorpus(argsMap cli.FlagMap) {

	_, err := client.client.TtsApi.DeleteCorpus(client.ctx, *argsMap["project_id"].Value,
		*argsMap["corpus_name"].Value)

	if nil == err {

		client.output.printInfo("Corpus deleted: " + *argsMap["corpus_name"].Value + " in project: " +
			*argsMap["project_id"].Value)

	} else {

		client.output.printError("DeleteCorpus", err)
	}
}

// Get delete corpus input attributes
func (client *TtsClient) GetDeleteCorpusProperties() map[string]bool {

	flagMap := map[string]bool{"project_id": true, "corpus_name": true}

	return flagMap
}
