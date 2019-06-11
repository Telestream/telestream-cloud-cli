package telestream

import (
	"context"
	"fmt"

	"tcs-cli/cli"
	"github.com/Telestream/telestream-cloud-go-sdk/flip"
)

// client holds output writer, flip configuration, flip client and context
type FlipClient struct {
	config *flip.Configuration
	client *flip.APIClient
	ctx    context.Context
	output ServiceOutput
}

// Creates new telestream flip client with given output writer and X API key
func NewFlipClient(xApiKey string, addHearderKey string, addHeaderVal string, output ServiceOutput) *FlipClient {

	client := new(FlipClient)
	client.config = flip.NewConfiguration()

	if addHearderKey != "" && addHeaderVal != "" {

		client.config.AddDefaultHeader(addHearderKey, addHeaderVal)
	}

	client.client = flip.NewAPIClient(client.config)
	client.ctx = context.WithValue(context.Background(), flip.ContextAPIKey, flip.APIKey{Key: xApiKey})
	client.output = output

	return client
}

// List all factories to output
func (client *FlipClient) ListFactories(argsMap cli.FlagMap) {

	opts, pageErr := getPageOpt(&argsMap)
	if pageErr != nil {

		client.output.printError("ListFactories", pageErr)
		return
	}

	factoriesCollection, _, err := client.client.FlipApi.Factories(client.ctx, opts)

	storageMap := map[int32]string{0: "S3", 1: "Google Cloud Storage", 2: "FTP storage", 5: "Flip storage",
		8: "FASP storage", 9: "Azure Blob Storage"}

	colNames := []interface{}{"NAME", "ID", "CREATED_AT", "STORE_ID", "OUTPUT_PATH_FORMAT"}
	rows := [][]interface{}{}

	if nil == err {

		for _, factory := range factoriesCollection.Factories {
			rows = append(rows, []interface{}{factory.Name, factory.Id, factory.CreatedAt,
				storageMap[factory.StorageProvider], factory.OutputsPathFormat})
		}

		client.output.printTable(colNames, rows)
	} else {

		client.output.printError("ListFactories", err)
	}
}

// Get list factories input attributes
func (client *FlipClient) GetListFactoriesProperties() map[string]bool {

	flagMap := map[string]bool{}
	addPageOpt(flagMap)

	return flagMap
}

// Print factory description on output
func (client *FlipClient) DescribeFactory(argsMap cli.FlagMap) {

	factoryDesc, _, err := client.client.FlipApi.Factory(client.ctx, *argsMap["factory_id"].Value,
		map[string]interface{}{})

	if nil == err {

		client.output.printStructContent(&factoryDesc)

	} else {

		client.output.printError("DescribeFatories", err)
	}
}

// Get describe factory input attributes
func (client *FlipClient) GetDescribeFactoryProperties() map[string]bool {

	flagMap := map[string]bool{"factory_id": true}

	return flagMap
}

// List all profiles for given factory on output
func (client *FlipClient) ListProfiles(argsMap cli.FlagMap) {

	opts, pageErr := getPageOpt(&argsMap)
	if pageErr != nil {

		client.output.printError("ListProfiles", pageErr)
		return
	}

	profilesCollection, _, err := client.client.FlipApi.Profiles(client.ctx, *argsMap["factory_id"].Value,
		opts)

	colNames := []interface{}{"NAME", " ID", "CREATED_AT", "FORMATS", "SIZE", "VIDEO_BITRATE", "AUDIO_BITRATE"}
	rows := [][]interface{}{}

	if nil == err {

		for _, profile := range profilesCollection.Profiles {

			vBitrate := fmt.Sprint(profile.VideoBitrate)
			aBitrate := fmt.Sprint(profile.AudioBitrate)

			if "0" == vBitrate {
				vBitrate = "auto"
			}

			if "0" == aBitrate {
				aBitrate = "auto"
			}

			rows = append(rows, []interface{}{profile.Name, profile.Id, profile.CreatedAt, profile.Title + " " + profile.AudioCodec,
				fmt.Sprint(profile.Width) + "x" + fmt.Sprint(profile.Height), aBitrate, vBitrate})
		}

		client.output.printTable(colNames, rows)

	} else {

		client.output.printError("ListProfiles", err)
	}
}

// Get list profile input attributes
func (client *FlipClient) GetListProfilesProperties() map[string]bool {

	flagMap := map[string]bool{"factory_id": true}
	addPageOpt(flagMap)

	return flagMap
}

// Print profile given by factory_id and profie_id/profile_name description
func (client *FlipClient) DescribeProfile(argsMap cli.FlagMap) {

	id_or_name := ""

	if val, ok := argsMap["profile_id"]; ok && *val.Value != "" {

		id_or_name = *val.Value
	} else if val, ok := argsMap["profile_name"]; ok && *val.Value != "" {

		id_or_name = *val.Value
	}

	if id_or_name != "" {

		profile, _, err := client.client.FlipApi.Profile(client.ctx, id_or_name, *argsMap["factory_id"].Value,
			map[string]interface{}{})

		if nil == err {

			client.output.printStructContent(&profile)

		} else {

			client.output.printError("DescribeProfile", err)
		}

	} else {

		client.output.printInfo("DescribeProfile: no profile_id or profile_name")
	}
}

// Get describe profile by name input attributes
func (client *FlipClient) GetDescribeProfileProperties() map[string]bool {

	flagMap := map[string]bool{"factory_id": true, "profile_name": false, "profile_id": true}

	return flagMap
}

// Create new profile in factory (selected by factory_id), print new profile description on output
func (client *FlipClient) CreateProfile(argsMap cli.FlagMap) {

	factory_id := *argsMap["factory_id"].Value
	delete(argsMap, "factory_id")

	newProfile := flip.ProfileBody{}
	propertiesToStruct(&newProfile, argsMap)

	profileDesc, _, err := client.client.FlipApi.CreateProfile(client.ctx, factory_id, newProfile,
		map[string]interface{}{})

	if nil == err {

		client.output.printStructContent(&profileDesc)

	} else {

		client.output.printError("CreateProfile", err)
	}
}

// Get new profile all input attributes
func (client *FlipClient) GetCreateProfileProperties() map[string]bool {

	flagMap := map[string]bool{"factory_id": true}
	profile := flip.ProfileBody{}

	jsonFields := structToProperties(&profile)

	for _, field := range jsonFields {

		flagMap[field] = false
	}

	flagMap["preset_name"] = true

	return flagMap
}

// Delete profile given by factory_id and profile_id, print result on output
func (client *FlipClient) DeleteProfile(argsMap cli.FlagMap) {

	factory_id := *argsMap["factory_id"].Value
	id := *argsMap["profile_id"].Value

	newProfile := flip.ProfileBody{}
	propertiesToStruct(&newProfile, argsMap)

	profileDel, _, err := client.client.FlipApi.DeleteProfile(client.ctx, id, factory_id)

	if nil == err {

		client.output.printStructContent(&profileDel)

	} else {

		client.output.printError("DeleteProfile", err)
	}
}

// Get delete attribute input arguments
func (client *FlipClient) GetDeleteProfileProperties() map[string]bool {

	flagMap := map[string]bool{"factory_id": true, "profile_id": true}

	return flagMap
}

// Update profile and print updated profile description
func (client *FlipClient) UpdateProfile(argsMap cli.FlagMap) {

	factory_id := *argsMap["factory_id"].Value
	id := *argsMap["profile_id"].Value
	delete(argsMap, "factory_id")
	delete(argsMap, "profile_id")

	newProfile := flip.ProfileBody{}
	propertiesToStruct(&newProfile, argsMap)

	profileDesc, _, err := client.client.FlipApi.UpdateProfile(client.ctx, id, factory_id, newProfile,
		map[string]interface{}{})

	if nil == err {

		client.output.printStructContent(&profileDesc)

	} else {

		client.output.printError("UpdateProfile", err)
	}
}

// Get update profile all input arguments
func (client *FlipClient) GetUpdateProfileProperties() map[string]bool {

	flagMap := map[string]bool{"factory_id": true, "profile_id": true}
	profile := flip.ProfileBody{}

	jsonFields := structToProperties(&profile)

	for _, field := range jsonFields {

		flagMap[field] = false
	}

	return flagMap
}

// List all videos in factory (given by factory_id) on output
func (client *FlipClient) ListVideos(argsMap cli.FlagMap) {

	opts, pageErr := getPageOpt(&argsMap)
	if pageErr != nil {

		client.output.printError("ListVideos", pageErr)
		return
	}

	videosCollection, _, err := client.client.FlipApi.Videos(client.ctx, *argsMap["factory_id"].Value,
		opts)

	colNames := []interface{}{"ORIGINAL_NAME", " ID", "CREATED_AT", "STATUS", "VIDEO_BITRATE", "AUDIO_BITRATE"}
	rows := [][]interface{}{}

	if nil == err {

		for _, video := range videosCollection.Videos {

			vBitrate := fmt.Sprint(video.VideoBitrate)
			aBitrate := fmt.Sprint(video.AudioBitrate)

			if "0" == vBitrate {
				vBitrate = "auto"
			}

			if "0" == aBitrate {
				aBitrate = "auto"
			}

			rows = append(rows, []interface{}{video.OriginalFilename, video.Id, video.CreatedAt,
				video.Status, aBitrate, aBitrate})
		}

		client.output.printTable(colNames, rows)
	} else {

		client.output.printError("ListVideos", err)
	}
}

// Get describe video all input attributes
func (client *FlipClient) GetListVideosProperties() map[string]bool {

	flagMap := map[string]bool{"factory_id": true}
	addPageOpt(flagMap)

	return flagMap
}

// Print video given by factory_id and video_id on output
func (client *FlipClient) DescribeVideo(argsMap cli.FlagMap) {

	video, _, err := client.client.FlipApi.Video(client.ctx, *argsMap["video_id"].Value, *argsMap["factory_id"].Value)

	encodingsCollection, _, errEnc := client.client.FlipApi.Encodings(client.ctx, *argsMap["factory_id"].Value,
		map[string]interface{}{"videoId": *argsMap["video_id"].Value})

	if nil == err && nil == errEnc {

		client.output.printStructContent(&video)

		encodingIds := ""
		for _, encoding := range encodingsCollection.Encodings {
			encodingIds += encoding.Id + ", "
		}

		client.output.printInfo("Encoding ids: " + encodingIds)

	} else {

		client.output.printError("DescribeVideo", err)
	}
}

// Get describe video all input attributes
func (client *FlipClient) GetDescribeVideoProperties() map[string]bool {

	flagMap := map[string]bool{"factory_id": true, "video_id": true}

	return flagMap
}

// Create new video and print new video description on output
func (client *FlipClient) CreateVideo(argsMap cli.FlagMap) {

	factory_id := *argsMap["factory_id"].Value
	delete(argsMap, "factory_id")

	newVideo := flip.CreateVideoBody{}
	propertiesToStruct(&newVideo, argsMap)

	videoDesc, _, err := client.client.FlipApi.CreateVideo(client.ctx, factory_id, newVideo)

	if nil == err {

		client.output.printStructContent(&videoDesc)

	} else {

		client.output.printError("CreateVideo", err)
	}
}

// Get create video all input attributes
func (client *FlipClient) GetCreateVideoProperties() map[string]bool {

	flagMap := map[string]bool{"factory_id": true}
	video := flip.CreateVideoBody{}

	jsonFields := structToProperties(&video)

	for _, field := range jsonFields {

		flagMap[field] = false
	}

	flagMap["source_url"] = true

	return flagMap
}

// Cancel video given by factory_id and video_id, print result on output
func (client *FlipClient) CancelVideo(argsMap cli.FlagMap) {

	factory_id := *argsMap["factory_id"].Value
	id := *argsMap["video_id"].Value

	videoCancel, _, err := client.client.FlipApi.CancelVideo(client.ctx, id, factory_id)

	if nil == err {

		client.output.printStructContent(&videoCancel)

	} else {

		client.output.printError("CancelVideo", err)
	}
}

// Get cancel video input attributes
func (client *FlipClient) GetCancelVideoProperties() map[string]bool {

	flagMap := map[string]bool{"factory_id": true, "video_id": true}

	return flagMap
}

// Delete video given by factory_id and video_id and print result
func (client *FlipClient) DeleteVideo(argsMap cli.FlagMap) {

	factory_id := *argsMap["factory_id"].Value
	id := *argsMap["video_id"].Value

	videoDelete, _, err := client.client.FlipApi.DeleteVideo(client.ctx, id, factory_id)

	if nil == err {

		client.output.printStructContent(&videoDelete)

	} else {

		client.output.printError("DeleteVideo", err)
	}
}

// Get delete video input attributes
func (client *FlipClient) GetDeleteVideoProperties() map[string]bool {

	flagMap := map[string]bool{"factory_id": true, "video_id": true}

	return flagMap
}

// List encodings for given factory_id (and video_id - optional)
func (client *FlipClient) ListEncodings(argsMap cli.FlagMap) {

	opts, pageErr := getPageOpt(&argsMap)
	if pageErr != nil {

		client.output.printError("ListEncodings", pageErr)
		return
	}

	factory_id := *argsMap["factory_id"].Value
	delete(argsMap, "factory_id")

	if *argsMap["video_id"].Value != "" {

		opts["videoId"] = *argsMap["video_id"].Value
	}

	encodingsCollection, _, err := client.client.FlipApi.Encodings(client.ctx, factory_id, opts)

	colNames := []interface{}{"ID", "CREATED_AT", "STATUS", "FILE_SIZE", "VIDEO_ID"}
	rows := [][]interface{}{}

	if nil == err {

		for _, encoding := range encodingsCollection.Encodings {

			rows = append(rows, []interface{}{encoding.Id, encoding.CreatedAt, encoding.Status, fmt.Sprint(encoding.FileSize), encoding.VideoId})
		}

		client.output.printTable(colNames, rows)
	} else {

		client.output.printError("ListEncodings", err)
	}
}

// Get list encodings inout attributes
func (client *FlipClient) GetListEncodingsProperties() map[string]bool {

	flagMap := map[string]bool{"factory_id": true, "video_id": false}
	addPageOpt(flagMap)

	return flagMap
}

// Print encoding description given by factory_id and encoding_id
func (client *FlipClient) DescribeEncoding(argsMap cli.FlagMap) {

	encoding, _, err := client.client.FlipApi.Encoding(client.ctx, *argsMap["encoding_id"].Value,
		*argsMap["factory_id"].Value, map[string]interface{}{})

	if nil == err {

		client.output.printStructContent(&encoding)
	} else {

		client.output.printError("DescribeVideos", err)
	}
}

// Get describe encoding input attributes
func (client *FlipClient) GetDescribeEncodingProperties() map[string]bool {

	flagMap := map[string]bool{"factory_id": true, "encoding_id": true}

	return flagMap
}

// Delete encoding given by factory_id an encoding_id, print result on output
func (client *FlipClient) DeleteEncoding(argsMap cli.FlagMap) {

	deleteEncoding, _, err := client.client.FlipApi.DeleteEncoding(client.ctx, *argsMap["encoding_id"].Value,
		*argsMap["factory_id"].Value)

	if nil == err {

		client.output.printStructContent(&deleteEncoding)
	} else {

		client.output.printError("DeleteEncoding", err)
	}
}

// Get delete encoding input attributes
func (client *FlipClient) GetDeleteEncodingProperties() map[string]bool {

	flagMap := map[string]bool{"factory_id": true, "encoding_id": true}

	return flagMap
}

// Print signed urls for specific encoding given by factory_id and encoding_id
func (client *FlipClient) SignedUrlsEncoding(argsMap cli.FlagMap) {

	EncodingSignedUrls, _, err := client.client.FlipApi.SignedEncodingUrls(client.ctx, *argsMap["encoding_id"].Value,
		*argsMap["factory_id"].Value)

	if nil == err {

		client.output.printStructContent(&EncodingSignedUrls)
	} else {

		client.output.printError("SignedUrlsEncoding", err)
	}
}

// Get signed urls encoding input attributes
func (client *FlipClient) GetSignedUrlsEncodingProperties() map[string]bool {

	flagMap := map[string]bool{"factory_id": true, "encoding_id": true}

	return flagMap
}

// Cancel encoding specified by factory_id and encoding_id, print result on output
func (client *FlipClient) CancelEncoding(argsMap cli.FlagMap) {

	cancelEncoding, _, err := client.client.FlipApi.CancelEncoding(client.ctx, *argsMap["encoding_id"].Value,
		*argsMap["factory_id"].Value)

	if nil == err {

		client.output.printStructContent(&cancelEncoding)
	} else {

		client.output.printError("CancelEncoding", err)
	}
}

// Get cancel encoding input attributes
func (client *FlipClient) GetCancelEncodingProperties() map[string]bool {

	flagMap := map[string]bool{"factory_id": true, "encoding_id": true}

	return flagMap
}
