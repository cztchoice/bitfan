package api

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	fqdn "github.com/ShowMax/go-fqdn"
	"github.com/gin-gonic/gin"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/vjeantet/bitfan/core"
	"github.com/vjeantet/bitfan/lib"
)

func getPipelineAssets(c *gin.Context) {
	var assets []Asset
	var id string

	id = c.Params.ByName("id")

	assets = []Asset{}
	ppls := core.Pipelines()
	for _, p := range ppls {
		if p.Uuid == id {

			loc, _ := lib.NewLocation(p.ConfigLocation, core.DataLocation())
			for path, b64Content := range loc.AssetsContent() {
				assets = append(assets, Asset{Path: path, Content: b64Content})
			}
			//build assets
			c.JSON(200, assets)
			return
		}
	}

	c.JSON(404, gin.H{"error": "no pipelines(s) running"})

	// curl -i http://localhost:8080/api/v1/pipelines/xxx/assets
}

func getPipelines(c *gin.Context) {

	var pipelines map[string]Pipeline
	var err error

	pipelines = map[string]Pipeline{}
	ppls := core.Pipelines()
	for _, p := range ppls {
		pipelines[p.Uuid] = Pipeline{
			Uuid:               p.Uuid,
			StartedAt:          p.StartedAt,
			Label:              p.Label,
			ConfigLocation:     p.ConfigLocation,
			ConfigHostLocation: p.ConfigHostLocation,
		}
	}

	if err == nil {
		c.JSON(200, pipelines)
	} else {
		c.JSON(404, gin.H{"error": "no pipelines(s) running"})
	}
	// curl -i http://localhost:8080/api/v1/pipelines
}

func getPipeline(c *gin.Context) {
	var pipeline Pipeline
	var id string

	id = c.Params.ByName("id")

	pipeline = Pipeline{}
	ppls := core.Pipelines()
	for _, p := range ppls {
		if p.Uuid == id {
			pipeline = Pipeline{
				Uuid:               p.Uuid,
				StartedAt:          p.StartedAt,
				Label:              p.Label,
				ConfigLocation:     p.ConfigLocation,
				ConfigHostLocation: p.ConfigHostLocation,
			}
			c.JSON(200, pipeline)
			return
		}
	}

	c.JSON(404, gin.H{"error": "no pipelines(s) running"})

}
func b64Decode(code string, dest string) error {
	buff, err := base64.StdEncoding.DecodeString(code)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(dest, buff, 07440); err != nil {
		return err
	}

	return nil
}
func addPipeline(c *gin.Context) {

	// ID, err := core.StartPipeline(&starter.Pipeline, starter.Agents)

	// b, err := ioutil.ReadAll(c.Request.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// pp.Println("c-->", string(b))

	//Bind request data
	var pipeline Pipeline
	err := c.BindJSON(&pipeline)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var cwd string
	if fqdn.Get() != pipeline.ConfigHostLocation { // Copy Assets
		// save Assets
		// directory = $data / remote / UUID /
		uidString := pipeline.Uuid
		if uidString == "" {
			uid, _ := uuid.NewV4()
			uidString = uid.String()
		}

		uidString = fmt.Sprintf("%s_%d", uidString, time.Now().Unix())

		cwd = filepath.Join(core.DataLocation(), "_pipelines", uidString)
		core.Log().Debugf("configuration %s stored to %s", uidString, cwd)
		os.MkdirAll(cwd, os.ModePerm)

		//Save assets to cwd
		for _, asset := range pipeline.Assets {
			dest := filepath.Join(cwd, asset.Path)
			dir := filepath.Dir(dest)
			os.MkdirAll(dir, os.ModePerm)
			b64Decode(asset.Content, dest)
			core.Log().Debugf("configuration %s asset %s stored", uidString, asset.Path)
		}

		pipeline.ConfigLocation = filepath.Join(cwd, pipeline.ConfigLocation)
		core.Log().Debugf("configuration %s pipeline %s ready to be loaded", uidString, pipeline.ConfigLocation)
	}

	var loc *lib.Location

	if pipeline.Content != "" {
		loc, err = lib.NewLocationContent(pipeline.Content, cwd)
	} else {
		loc, err = lib.NewLocation(pipeline.ConfigLocation, cwd)
	}

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ppl := loc.ConfigPipeline()
	if pipeline.Label != "" {
		ppl.Name = pipeline.Label
	}

	if pipeline.Uuid != "" {
		ppl.Uuid = pipeline.Uuid
	}

	agt, err := loc.ConfigAgents()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	UUID, err := core.StartPipeline(&ppl, agt)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	core.Log().Debugf("Pipeline %s started", pipeline.Label)
	c.Redirect(http.StatusFound, fmt.Sprintf("/api/v1/pipelines/%s", UUID))
	return

	// c.JSON(200, gin.H{"statut": "started"})

	// curl -i -X DELETE http://localhost:8080/api/v1/pipelines/1
}

func deletePipeline(c *gin.Context) {

	var err error
	id := c.Params.ByName("id")
	if id == "" {
		c.JSON(500, gin.H{"error": fmt.Sprintf("malformed id : %s", c.Params.ByName("id"))})
		return
	}

	err = core.StopPipeline(id)

	if err == nil {
		c.JSON(200, gin.H{"id": c.Params.ByName("id"), "status": "deleted"})
		core.Log().Debugf("Pipeline %s stopped", id)
	} else {
		c.JSON(404, gin.H{"error": err.Error()})
	}

	// curl -i -X DELETE http://localhost:8080/api/v1/pipelines/1
}

func getDocs(c *gin.Context) {

	data := docsByKind("input")
	data = append(data, docsByKind("filter")...)
	data = append(data, docsByKind("output")...)

	c.JSON(200, data)
	// curl -i http://localhost:8080/api/v1/docs
}

func getDocsInputs(c *gin.Context) {
	c.JSON(200, docsByKind("input"))
	// curl -i http://localhost:8080/api/v1/docs/inputs
}

func getDocsFilters(c *gin.Context) {
	c.JSON(200, docsByKind("filter"))
	// curl -i http://localhost:8080/api/v1/docs/filters
}

func getDocsOutputs(c *gin.Context) {
	c.JSON(200, docsByKind("output"))
	// curl -i http://localhost:8080/api/v1/docs/outputs
}

func getDocsInputsByName(c *gin.Context) {
	data, err := docsByKindByName("input", c.Params.ByName("name"))
	if err != nil {
		c.JSON(404, gin.H{"error": "not found"})
	} else {
		c.JSON(200, data)
	}
}

func getDocsFiltersByName(c *gin.Context) {
	data, err := docsByKindByName("filter", c.Params.ByName("name"))
	if err != nil {
		c.JSON(404, gin.H{"error": "not found"})
	} else {
		c.JSON(200, data)
	}
}

func getDocsOutputsByName(c *gin.Context) {
	data, err := docsByKindByName("output", c.Params.ByName("name"))
	if err != nil {
		c.JSON(404, gin.H{"error": "not found"})
	} else {
		c.JSON(200, data)
	}
}