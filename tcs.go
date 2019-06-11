package main

import (
	"fmt"
	"os"
	"tcs-cli/cli"
	"tcs-cli/telestream"
	"gopkg.in/ini.v1"
)

const configFileName = ".tcs-credentials"

func createTtsProjectsCommands(client *telestream.TtsClient) []cli.CommandBaseInterface {

	// projects list command
	projectsListCmd := cli.NewCommand("list", client.ListProjects, "List all Projects")

	// factories describe command
	projectsDescribeCmd := cli.NewFlaggedCommand("describe", "project_id", client.DescribeProject,
		client.GetDescribeProjectProperties(), "describes project by project_id")

	// projects create command
	projectsCreateCmd := cli.NewFlaggedCommand("create", "", client.CreateProject,
		client.GetCreateProjectProperties(), "creates project")

	// projects delete command
	projectsDeleteCmd := cli.NewFlaggedCommand("delete", "project_id", client.DeleteProject,
		client.GetDeleteProjectProperties(), "deletes project")

	// projects update command
	projectsUpdateCmd := cli.NewFlaggedCommand("update", "project_id", client.UpdateProject,
		client.GetUpdateProjectProperties(), "updates project")

	return []cli.CommandBaseInterface{projectsListCmd, projectsDescribeCmd, projectsCreateCmd,
		projectsDeleteCmd, projectsUpdateCmd}
}

func createTtsJobsCommands(client *telestream.TtsClient) []cli.CommandBaseInterface {

	// jobs list command
	jobsListCmd := cli.NewFlaggedCommand("list", "", client.ListJobs,
		client.GetListJobsProperties(), "List all jobs")

	// jobs create command
	jobsCreateCmd := cli.NewFlaggedCommand("create", "", client.CreateJob,
		client.GetCreateJobProperties(), "Create job")

	// jobs describe command
	jobsDescribeCmd := cli.NewFlaggedCommand("describe", "job_id", client.DescribeJob,
		client.GetDescribeJobProperties(), "Describe job")

	// jobs delete command
	jobsDeleteCmd := cli.NewFlaggedCommand("delete", "job_id", client.DeleteJob,
		client.GetDeleteJobProperties(), "Delete job")

	// jobs result command
	jobsJobResultCmd := cli.NewFlaggedCommand("result", "job_id", client.JobResult,
		client.GetJobResultProperties(), "Describe job result")

	// jobs outputs command
	jobsJobOutputsCmd := cli.NewFlaggedCommand("outputs", "job_id", client.JobOutputs,
		client.GetJobOutputsProperties(), "Describe job outputs")

	return []cli.CommandBaseInterface{jobsListCmd, jobsCreateCmd, jobsDescribeCmd, jobsDeleteCmd,
		jobsJobResultCmd, jobsJobOutputsCmd}
}

func createTtsCorporaCommands(client *telestream.TtsClient) []cli.CommandBaseInterface {

	// corpora list command
	corporaListCmd := cli.NewFlaggedCommand("list", "", client.ListCorpora, client.GetListCorporaProperties(),
		"List all corpora")

	// corpora describe command
	corporaDescribeCmd := cli.NewFlaggedCommand("describe", "corpus_name", client.DescribeCorpus,
		client.GetDescribeCorpusProperties(), "describes corpus by project_id and corpus name")

	// corpora create command
	corporaCreateCmd := cli.NewFlaggedCommand("create", "", client.CreateCorpus,
		client.GetCreateCorpusProperties(), "creates corpus")

	// corpora delete command
	corporaDeleteCmd := cli.NewFlaggedCommand("delete", "corpus_name", client.DeleteCorpus,
		client.GetDeleteCorpusProperties(), "deletes corpus")

	return []cli.CommandBaseInterface{corporaListCmd, corporaDescribeCmd, corporaCreateCmd,
		corporaDeleteCmd}
}

func createTtsCommands(client *telestream.TtsClient) []cli.CommandBaseInterface {

	// projects subcommand
	projectsCmd := cli.NewSubCommand("projects", createTtsProjectsCommands(client),
		"manage Telesteam cloud tts service projects")

	// jobs subcommand
	jobsCmd := cli.NewSubCommand("jobs", createTtsJobsCommands(client),
		"manage Telesteam cloud tts service jobs")

	// corpora subcommand
	corporaCmd := cli.NewSubCommand("corpora", createTtsCorporaCommands(client),
		"manage Telestream cloud tts service corpora")

	// tts subcommand
	flipCmd := cli.NewSubCommand("tts", []cli.CommandBaseInterface{projectsCmd, jobsCmd, corporaCmd},
		"manage your tts service")

	return []cli.CommandBaseInterface{flipCmd}
}

func createFactoriesCommands(client *telestream.FlipClient) []cli.CommandBaseInterface {

	// factories list command
	factoriesListCmd := cli.NewFlaggedCommand("list", "", client.ListFactories, client.GetListFactoriesProperties(),
		"List all factories")

	// factories describe command
	factoriesDescribeCmd := cli.NewFlaggedCommand("describe", "factory_id", client.DescribeFactory,
		client.GetDescribeFactoryProperties(), "describes factory by factory_id")

	return []cli.CommandBaseInterface{factoriesListCmd, factoriesDescribeCmd}
}

func createProfilesCommands(client *telestream.FlipClient) []cli.CommandBaseInterface {

	// profiles list command
	profilesListCmd := cli.NewFlaggedCommand("list", "", client.ListProfiles,
		client.GetListProfilesProperties(), "lists profiles by factory_id")

	// profiles describe command
	profilesDescribeByNameCmd := cli.NewFlaggedCommand("describe", "profile_id", client.DescribeProfile,
		client.GetDescribeProfileProperties(), "describes profile by factory_id and its name or id")

	// profiles create createcommand
	profilesCreateCmd := cli.NewFlaggedCommand("create", "", client.CreateProfile,
		client.GetCreateProfileProperties(), "creates profile")

	// profiles delete createcommand
	profilesDeleteCmd := cli.NewFlaggedCommand("delete", "profile_id", client.DeleteProfile,
		client.GetDeleteProfileProperties(), "deletes profile")

	// profiles update createcommand
	profilesUpdateCmd := cli.NewFlaggedCommand("update", "profile_id", client.UpdateProfile,
		client.GetUpdateProfileProperties(), "updates profile")

	return []cli.CommandBaseInterface{profilesListCmd, profilesDescribeByNameCmd,
		profilesCreateCmd, profilesDeleteCmd, profilesUpdateCmd}
}

func createVideosCommands(client *telestream.FlipClient) []cli.CommandBaseInterface {

	// videos list command
	videosListCmd := cli.NewFlaggedCommand("list", "", client.ListVideos,
		client.GetListVideosProperties(), "lists videos by factory_id")

	// videos describe command
	videosDescribeCmd := cli.NewFlaggedCommand("describe", "video_id", client.DescribeVideo,
		client.GetDescribeVideoProperties(), "describes videos by factory_id and its id")

	// videos create command
	videosCreateCmd := cli.NewFlaggedCommand("create", "", client.CreateVideo,
		client.GetCreateVideoProperties(), "creates video")

	// videos cancel command
	videosCancelCmd := cli.NewFlaggedCommand("cancel", "video_id", client.CancelVideo,
		client.GetCancelVideoProperties(), "cancels video")

	// videos delete command
	videosDeleteCmd := cli.NewFlaggedCommand("delete", "video_id", client.DeleteVideo,
		client.GetDeleteVideoProperties(), "deletes video")

	return []cli.CommandBaseInterface{videosListCmd, videosDescribeCmd, videosCreateCmd,
		videosCancelCmd, videosDeleteCmd}
}

func createEncodingsCommands(client *telestream.FlipClient) []cli.CommandBaseInterface {

	// encodings list command
	encodingsListCmd := cli.NewFlaggedCommand("list", "", client.ListEncodings,
		client.GetListEncodingsProperties(), "lists all encodings")

	// encodings describe command
	encodingsDescribeCmd := cli.NewFlaggedCommand("describe", "encoding_id", client.DescribeEncoding,
		client.GetDescribeEncodingProperties(), "describes encoding by factory_id and its id")

	// encodings cancel command
	cancelDescribeCmd := cli.NewFlaggedCommand("cancel", "encoding_id", client.CancelEncoding,
		client.GetCancelEncodingProperties(), "cancels encoding by factory_id and encoding id")

	// encodings signed-urls command
	signedUrlsDescribeCmd := cli.NewFlaggedCommand("signed-urls", "encoding_id", client.SignedUrlsEncoding,
		client.GetSignedUrlsEncodingProperties(), "signed-urls by factory_id and encoding id")

	// encodings delete command
	deleteDescribeCmd := cli.NewFlaggedCommand("delete", "encoding_id", client.DeleteEncoding,
		client.GetDeleteEncodingProperties(), "deletes encoding by factory_id and encoding id")

	return []cli.CommandBaseInterface{encodingsListCmd, encodingsDescribeCmd, cancelDescribeCmd,
		signedUrlsDescribeCmd, deleteDescribeCmd}
}

func createFlipCommands(client *telestream.FlipClient) []cli.CommandBaseInterface {

	// factories subcommand
	factoriesCmd := cli.NewSubCommand("factories", createFactoriesCommands(client),
		"manage Telesteam cloud service factories")

	// profiles subcommand
	profilesCmd := cli.NewSubCommand("profiles", createProfilesCommands(client),
		"manage Telesteam cloud service profiles")

	// videos subcommand
	videosCmd := cli.NewSubCommand("videos", createVideosCommands(client),
		"menage Telestream cloud service videos")

	// encodings subcommand
	encodingsCmd := cli.NewSubCommand("encodings", createEncodingsCommands(client),
		"menage Telestream cloud service encodings")

	// flip subcommand
	flipCmd := cli.NewSubCommand("flip", []cli.CommandBaseInterface{factoriesCmd, profilesCmd, videosCmd, encodingsCmd},
		"manage your flip service")

	return []cli.CommandBaseInterface{flipCmd}
}

func createConfig(argsMap cli.FlagMap) {

	key := *argsMap["api_key"].Value

	var home string
	var err error

	if home, err = os.UserHomeDir(); err != nil {

		fmt.Println("Cannot get home directory: ", err.Error())
		return
	}

	configFilePath := home + "/" + configFileName
	os.Remove(configFilePath)

	cfg := ini.Empty()
	cfg.Section("default").Key("api_key").SetValue(key)

	if _, err = os.Create(configFilePath); err != nil {

		fmt.Println("Cannot create configuration file: ", err.Error())
		return
	}

	if err = cfg.SaveTo(configFilePath); err != nil {

		fmt.Println("Cannot save configuration file: ", err.Error())
		return
	}

	fmt.Println("Credentials saved")
}

func readConfig() string {

	var home string
	var err error

	if home, err = os.UserHomeDir(); err != nil {

		fmt.Println("Cannot get home directory: ", err.Error())
		return ""
	}

	cfg, err := ini.Load(home + "/" + configFileName)
	if err != nil {

		return ""
	}

	return cfg.Section("default").Key("api_key").String()
}

func createTcsCli() {

	configureCmdStr := "configure"

	apiKey := readConfig()

	if apiKey == "" && len(os.Args) > 1 && os.Args[1] != configureCmdStr {

		fmt.Println("Firstly you should configure credentials")
		return
	}

	argvOutput := os.Args

	var flags map[string]*string

	argvOutput, flags = cli.GetAdditionalFlags(os.Args, map[string]string{"header_key": "additive http header key",
		"header_val": "additive http header val"})

	additionalHeaderKey := *flags["header_key"]
	additionalHeaderVal := *flags["header_val"]

	flipClient := telestream.NewFlipClient(apiKey, additionalHeaderKey, additionalHeaderVal, telestream.ServiceToStdOut)
	ttsClient := telestream.NewTtsClient(apiKey, additionalHeaderKey, additionalHeaderVal, telestream.ServiceToStdOut)

	// configure command
	configureCmd := cli.NewFlaggedCommand(configureCmdStr, "api_key", createConfig, map[string]bool{"api_key": true},
		"create configuration file for tsc command line tool with credentials that are used to interact with telestream cloud API")

	commands := createFlipCommands(flipClient)
	commands = append(commands, createTtsCommands(ttsClient)[0])
	commands = append(commands, configureCmd)

	cmdHndl := cli.NewCommandHandler("tcs", commands, map[string]string{"header_key": "additive http header key",
		"header_val": "additive http header value"})

	cmdHndl.ParseArgs(argvOutput)
}
