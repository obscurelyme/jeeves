package pom

import (
	"encoding/xml"
)

// Official Pom Schema https://maven.apache.org/xsd/maven-4.0.0.xsd
type Project struct {
	Attrs                  []xml.Attr              `xml:",any,attr"`
	XMLName                xml.Name                `xml:"project"`
	ModelVersion           string                  `xml:"modelVersion,omitempty"`
	Parent                 *Parent                 `xml:"parent,omitempty"`
	GroupId                string                  `xml:"groupId,omitempty"`
	ArtifactId             string                  `xml:"artifactId,omitempty"`
	Version                string                  `xml:"version,omitempty"`
	Packaging              string                  `xml:"packaging,omitempty"`
	Name                   string                  `xml:"name,omitempty"`
	Description            string                  `xml:"description,omitempty"`
	Url                    string                  `xml:"url,omitempty"`
	InceptionYear          string                  `xml:"inceptionYear,omitempty"`
	Organization           *Organization           `xml:"organization,omitempty"`
	Licenses               *Licenses               `xml:"licenses,omitempty"`
	Developers             *Developers             `xml:"developers,omitempty"`
	Contributors           *Contributors           `xml:"contributors,omitempty"`
	MailingLists           *MailingLists           `xml:"mailingLists,omitempty"`
	Prerequisites          *Prerequisites          `xml:"prerequisites,omitempty"`
	Modules                *Modules                `xml:"modules,omitempty"`
	Scm                    *Scm                    `xml:"scm,omitempty"`
	IssueManagement        *IssueManagement        `xml:"issueManagement,omitempty"`
	CiManagement           *CiManagement           `xml:"ciManagement,omitempty"`
	DistributionManagement *DistributionManagement `xml:"distributionManagement,omitempty"`
	Properties             *Any                    `xml:"properties,omitempty"`
	DependencyManagement   *DependencyManagement   `xml:"dependencyManagement,omitempty"`
	Dependencies           *Dependencies           `xml:"dependencies,omitempty"`
	Repositories           *Repositories           `xml:"repositories,omitempty"`
	PluginRepositories     *PluginRepositories     `xml:"pluginRepositories,omitempty"`
	Build                  *Build                  `xml:"build,omitempty"`
	// Deprecated: Now ignored by Maven.
	Reports   *Reports   `xml:"reports,omitempty"`
	Reporting *Reporting `xml:"reporting,omitempty"`
	Plugins   *Plugins   `xml:"plugins,omitempty"`
	// A listing of project-local build profiles which will modify the build process when activated.
	Profiles *Profiles `xml:"profiles,omitempty"`
}

type Prerequisites struct {
	Comment string `xml:",comment"`
	// For a plugin project (packaging is <code>maven-plugin</code>), the minimum version of Maven required to use the resulting plugin.
	Maven string `xml:"maven,omitempty"`
}

type Modules struct {
	Comment string   `xml:",comment"`
	Module  []string `xml:"module,omitempty"`
}

type Licenses struct {
	Comment string    `xml:",comment"`
	License []License `xml:"license,omitempty"`
}

type License struct {
	Comment      string `xml:",comment"`
	Name         string `xml:"name,omitempty"`
	Url          string `xml:"url,omitempty"`
	Distribution string `xml:"distribution,omitempty"`
	Comments     string `xml:"comments,omitempty"`
}

type CiManagement struct {
	Comment   string     `xml:",comment"`
	System    string     `xml:"system,omitempty"`
	Url       string     `xml:"url,omitempty"`
	Notifiers *Notifiers `xml:"notifiers,omitempty"`
}

type Notifiers struct {
	Comment  string     `xml:",comment"`
	Notifier []Notifier `xml:"notifier,omitempty"`
}

type Notifier struct {
	Comment       string `xml:",comment"`
	Type          string `xml:"type,omitempty"`
	SendOnError   bool   `xml:"sendOnError,omitempty"`
	SendOnFailure bool   `xml:"sendOnFailure,omitempty"`
	SendOnSuccess bool   `xml:"sendOnSuccess,omitempty"`
	SendOnWarning bool   `xml:"sendOnWarning,omitempty"`
	Address       string `xml:"address,omitempty"`
	Configuration *Any   `xml:"configuration,omitempty"`
}

type Scm struct {
	Comment string `xml:",comment"`
	/*
		The source control management system URL that describes the repository and how to connect to the repository.
		For more information, see the [URL format] and [list of supported SCMs]. This connection is read-only.

		[URL format]: https://maven.apache.org/scm/scm-url-format.html
		[list of supported SCMs]: https://maven.apache.org/scm/scms-overview.html

	*/
	Connection          string `xml:"connection,omitempty"`
	DeveloperConnection string `xml:"developerConnection,omitempty"`
	Tag                 string `xml:"tag,omitempty"`
	Url                 string `xml:"url,omitempty"`
}

type IssueManagement struct {
	Comment string `xml:",comment"`
	System  string `xml:"system,omitempty"`
	Url     string `xml:"url,omitempty"`
}

type DependencyManagement struct {
	Comment      string        `xml:",comment"`
	Dependencies *Dependencies `xml:"dependencies,omitempty"`
}

type Dependency struct {
	GroupId    string      `xml:"groupId,omitempty"`
	ArtifactId string      `xml:"artifactId,omitempty"`
	Version    string      `xml:"version,omitempty"`
	Type       string      `xml:"type,omitempty"`
	Classifier string      `xml:"classifier,omitempty"`
	Scope      string      `xml:"scope,omitempty"`
	SystemPath string      `xml:"systemPath,omitempty"`
	Exclusions *Exclusions `xml:"exclusions,omitempty"`
	Optional   string      `xml:"optional,omitempty"`
}

type Exclusions struct {
	Comment   string      `xml:",comment"`
	Exclusion []Exclusion `xml:"exclusion,omitempty"`
}

type Exclusion struct {
	Comment    string `xml:",comment"`
	ArtifactId string `xml:"artifactId,omitempty"`
	GroupId    string `xml:"groupId,omitempty"`
}

type Parent struct {
	Comment      string `xml:",comment"`
	GroupId      string `xml:"groupId,omitempty"`
	ArtifactId   string `xml:"artifactId,omitempty"`
	Version      string `xml:"version,omitempty"`
	RelativePath string `xml:"relativePath,omitempty"`
}

type Developers struct {
	Comment   string      `xml:",comment"`
	Developer []Developer `xml:"developer,omitempty"`
}

type Developer struct {
	Id              string `xml:"id,omitempty"`
	Name            string `xml:"name,omitempty"`
	Email           string `xml:"email,omitempty"`
	Url             string `xml:"url,omitempty"`
	Organization    string `xml:"organization,omitempty"`
	OrganizationUrl string `xml:"organizationUrl,omitempty"`
	Roles           *Roles `xml:"roles,omitempty"`
	Timezone        string `xml:"timezone:omitempty"`
	Properties      *Any   `xml:"properties,omitempty"`
}

type Roles struct {
	Comment string   `xml:",comment"`
	Role    []string `xml:"role,omitempty"`
}

type MailingLists struct {
	Comment     string        `xml:",comment"`
	MailingList []MailingList `xml:"mailingList,omitempty"`
}

type MailingList struct {
	Comment       string         `xml:",comment"`
	Name          string         `xml:"name,omitempty"`
	Subscribe     string         `xml:"subscribe,omitempty"`
	Unsubscribe   string         `xml:"unsubscribe,omitempty"`
	Post          string         `xml:"post,omitempty"`
	Archive       string         `xml:"archive,omitempty"`
	OtherArchives []OtherArchive `xml:"otherArchives,omitempty"`
}

type OtherArchive struct {
	Comment      string `xml:",comment"`
	OtherArchive string `xml:"otherArchive,omitempty"`
}

type Contributors struct {
	Comment     string        `xml:",comment"`
	Contributor []Contributor `xml:"contributor,omitempty"`
}

type Contributor struct {
	Developer
}

type Organization struct {
	Comment string `xml:",comment"`
	Name    string `xml:"name,omitempty"`
	Url     string `xml:"url,omitempty"`
}

type Dependencies struct {
	Comment    string       `xml:",comment"`
	Dependency []Dependency `xml:"dependency,omitempty"`
}

type DistributionManagement struct {
	Comment            string                `xml:",comment"`
	Repository         *DeploymentRepository `xml:"repository,omitempty"`
	SnapshotRepository *DeploymentRepository `xml:"snapshotRepository,omitempty"`
	Site               *Site                 `xml:"site,omitempty"`
	DownloadUrl        string                `xml:"downloadUrl,omitempty"`
	Reloction          *Relocation           `xml:"relocation,omitempty"`
	Status             string                `xml:"status,omitempty"`
}

type DeploymentRepository struct {
	Comment       string            `xml:",comment"`
	UniqueVersion bool              `xml:"uniqueVersion,omitempty"`
	Releases      *RepositoryPolicy `xml:"releases,omitempty"`
	Snapshots     *RepositoryPolicy `xml:"snapshots,omitempty"`
	Id            string            `xml:"id,omitempty"`
	Name          string            `xml:"name,omitempty"`
	Url           string            `xml:"url,omitempty"`
	Layout        string            `xml:"layout,omitempty"`
}

type Repositories struct {
	Comment    string       `xml:",comment"`
	Repository []Repository `xml:"repository,omitempty"`
}

type PluginRepositories struct {
	Comment    string       `xml:",comment"`
	Repository []Repository `xml:"repository,omitempty"`
}

type Repository struct {
	Comment   string            `xml:",comment"`
	Releases  *RepositoryPolicy `xml:"releases,omitempty"`
	Snapshots *RepositoryPolicy `xml:"snapshots,omitempty"`
	Id        string            `xml:"id,omitempty"`
	Name      string            `xml:"name,omitempty"`
	Url       string            `xml:"url,omitempty"`
	Layout    string            `xml:"layout,omitempty"`
}

// Download policy
type RepositoryPolicy struct {
	Comment        string `xml:",comment"`
	Enabled        string `xml:"enabled,omitempty"`
	UpdatePolicy   string `xml:"updatPolicy,omitempty"`
	ChecksumPolicy string `xml:"checksumPolicy,omitempty"`
}

type Site struct {
	Comment string `xml:",comment"`
	Id      string `xml:"id,omitempty"`
	Name    string `xml:"name,omitempty"`
	Url     string `xml:"url,omitempty"`
}

type Relocation struct {
	Comment    string `xml:",comment"`
	GroupId    string `xml:"groupId,omitempty"`
	ArtifactId string `xml:"artifactId,omitempty"`
	Version    string `xml:"version,omitempty"`
	Message    string `xml:"message,omitempty"`
}

type Reports struct {
	Comment string   `xml:",comment"`
	Report  []string `xml:"report,omitempty"`
}

type Reporting struct {
	ExcludeDefaults string         `xml:"excludeDefaults,omitempty"`
	OutputDirectory string         `xml:"outputDirectory,omitempty"`
	Plugins         *ReportPlugins `xml:"plugins,omitempty"`
}

type ReportPlugins struct {
	Comment string         `xml:",comment"`
	Plugins []ReportPlugin `xml:"plugins,omitempty"`
}

type ReportPlugin struct {
	Comment    string `xml:",comment"`
	GroupId    string `xml:"groupId,omitempty"`
	ArtifactId string `xml:"artifactId,omitempty"`
	Version    string `xml:"version,omitempty"`
	// Multiple specifications of a set of reports, each having (possibly) different configuration.
	// This is the reporting parallel to an <code>execution</code> in the build.
	ReportSets    *ReportSets `xml:"reportSets,omitempty"`
	Inherited     string      `xml:"inherited,omitempty"`
	Configuration *Any        `xml:"configuration,omitempty"`
}

type ReportSets struct {
	Comment   string      `xml:",comment"`
	ReportSet []ReportSet `xml:"reportSet,omitempty"`
}

type ReportSet struct {
	Id            string   `xml:"id,omitempty"`
	Reports       *Reports `xml:"reports,omitempty"`
	Inherited     string   `xml:"inherited,omitempty"`
	Configuration *Any     `xml:"configuration,omitempty"`
}

type PluginManagement struct {
	Comment string   `xml:",comment"`
	Plugins *Plugins `xml:"plugins,omitempty"`
}

type Plugins struct {
	Comment string   `xml:",comment"`
	Plugin  []Plugin `xml:"plugin,omitempty"`
}

type Plugin struct {
	Comment      string        `xml:",comment"`
	GroupId      string        `xml:"groupId,omitempty"`
	ArtifactId   string        `xml:"artifactId,omitempty"`
	Version      string        `xml:"version,omitempty"`
	Extensions   bool          `xml:"extensions,omitempty"`
	Executions   *Executions   `xml:"executions,omitempty"`
	Dependencies *Dependencies `xml:"dependencies,omitempty"`
	// Deprecated: Not used by Maven. Use Goals within Execution instead
	Goals         *Goals `xml:"goals,omitempty"`
	Inherited     string `xml:"inherited,omitempty"`
	Configuration *Any   `xml:"configuration,omitempty"`
}

type Executions struct {
	Comment   string      `xml:",comment"`
	Execution []Execution `xml:"execution,omitempty"`
}

type Execution struct {
	Id    string `xml:"id,omitempty"`
	Phase string `xml:"phase,omitempty"`
	Goals *Goals `xml:"goals,omitempty"`
}

type Goals struct {
	Comment string   `xml:",comment"`
	Goal    []string `xml:"goal,omitempty"`
}

type Resources struct {
	Comment string `xml:",comment"`
	/*
		Describe the resource target path. The path is relative to the target/classes directory

		IE: ${project.build.outputDirectory}
	*/
	TargetPath string    `xml:"targetPath,omitempty"`
	Filtering  string    `xml:"filtering,omitempty"`
	Directory  string    `xml:"directory,omitempty"`
	Includes   *Includes `xml:"includes,omitempty"`
	Excludes   *Excludes `xml:"excludes,omitempty"`
}

type Includes struct {
	Comment string   `xml:",comment"`
	Include []string `xml:"include,omitempty"`
}

type Excludes struct {
	Comment string   `xml:",comment"`
	Exclude []string `xml:"exclude,omitempty"`
}

type Filters struct {
	Comment string `xml:",comment"`
}

type BuildBase struct {
	Comment          string            `xml:",comment"`
	DefaultGoal      string            `xml:"defaultGoal,omitempty"`
	Resources        *Resources        `xml:"resources,omitempty"`
	TestResources    *Resources        `xml:"testResources,omitempty"`
	Directory        string            `xml:"directory,omitempty"`
	FinalName        string            `xml:"finalName,omitempty"`
	Filters          *Filters          `xml:"filters,omitempty"`
	PluginManagement *PluginManagement `xml:"pluginManagement,omitempty"`
	Plugins          *Plugins          `xml:"plugins,omitempty"`
}

type Build struct {
	Comment               string      `xml:",comment"`
	SourceDirectory       string      `xml:"sourceDirectory,omitempty"`
	ScriptSourceDirectory string      `xml:"scriptSourceDirectory,omitempty"`
	TestSourceDirectory   string      `xml:"testSourceDirectory,omitempty"`
	OutputDirectory       string      `xml:"outputDirectory,omitempty"`
	TestOutputDirectory   string      `xml:"testOutputDirectory,omitempty"`
	Extensions            *Extensions `xml:"extensions,omitempty"`
	BuildBase
}

type Extensions struct {
	Comment   string      `xml:",comment"`
	Extension []Extension `xml:"extension,omitempty"`
}

type Extension struct {
	Comment    string `xml:",comment"`
	GroupId    string `xml:"groupId,omitempty"`
	ArtifactId string `xml:"artifactId,omitempty"`
	Version    string `xml:"version,omitempty"`
}

type Profiles struct {
	Comment string    `xml:",comment"`
	Profile []Profile `xml:"profile,omitempty"`
}

// Modifications to the build process which is activated based on environmental parameters or command line arguments.
type Profile struct {
	Comment                string                  `xml:",comment"`
	Id                     string                  `xml:"id,omitempty"`
	Activation             *Activation             `xml:"activation,omitempty"`
	Build                  *BuildBase              `xml:"build,omitempty"`
	Modules                *Modules                `xml:"modules,omitempty"`
	DistributionManagement *DistributionManagement `xml:"distributionManagement,omitempty"`
	Properties             *Any                    `xml:"properties,omitempty"`
	Dependencies           *Dependencies           `xml:"dependencies,omitempty"`
	Repositories           *Repositories           `xml:"repositories,omitempty"`
	PluginRepositories     *PluginRepositories     `xml:"pluginRepositories,omitempty"`
	// Deprecated: Not used by Maven
	Reports   *Reports   `xml:"reports,omitempty"`
	Reporting *Reporting `xml:"reporting,omitempty"`
}

type Activation struct {
	Comment         string              `xml:",comment"`
	ActiveByDefault bool                `xml:"activeByDefault,omitempty"`
	JDK             string              `xml:"jdk,omitempty"`
	OS              *ActivationOS       `xml:"os,omitempty"`
	Property        *ActivationProperty `xml:"property,omitempty"`
	File            *ActivationFile     `xml:"file,omitempty"`
}

type ActivationProperty struct {
	Comment string `xml:",comment"`
	Name    string `xml:"name,omitempty"`
	Value   string `xml:"value,omitempty"`
}

type ActivationOS struct {
	Comment string `xml:",comment"`
	Name    string `xml:"name,omitempty"`
	Family  string `xml:"family,omitempty"`
	Arch    string `xml:"arch,omitempty"`
	Version string `xml:"version,omitempty"`
}

type ActivationFile struct {
	Comment string `xml:",comment"`
	// The name of the file that must be missing to activate the profile.
	Missing string `xml:"missing,omitempty"`
	// The name of the file that must exist to activate the profile.
	Exists string `xml:"exists,omitempty"`
}

type Any struct {
	XMLName     xml.Name
	Attrs       []xml.Attr `xml:"-"`
	Value       string     `xml:",chardata"`
	AnyElements []Any      `xml:",any"`
}
