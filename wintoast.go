package wintoast

import (
	"bytes"
	"os/exec"
	"syscall"
	"text/template"
)

type Audio string
type Duration string
type ActivationType string
type Scenario string

var toastTemplate *template.Template

const (
	Default        Audio = "ms-winsoundevent:Notification.Default"
	IM             Audio = "ms-winsoundevent:Notification.IM"
	Mail           Audio = "ms-winsoundevent:Notification.Mail"
	Reminder       Audio = "ms-winsoundevent:Notification.Reminder"
	SMS            Audio = "ms-winsoundevent:Notification.SMS"
	LoopingAlarm   Audio = "ms-winsoundevent:Notification.Looping.Alarm"
	LoopingAlarm2  Audio = "ms-winsoundevent:Notification.Looping.Alarm2"
	LoopingAlarm3  Audio = "ms-winsoundevent:Notification.Looping.Alarm3"
	LoopingAlarm4  Audio = "ms-winsoundevent:Notification.Looping.Alarm4"
	LoopingAlarm5  Audio = "ms-winsoundevent:Notification.Looping.Alarm5"
	LoopingAlarm6  Audio = "ms-winsoundevent:Notification.Looping.Alarm6"
	LoopingAlarm7  Audio = "ms-winsoundevent:Notification.Looping.Alarm7"
	LoopingAlarm8  Audio = "ms-winsoundevent:Notification.Looping.Alarm8"
	LoopingAlarm9  Audio = "ms-winsoundevent:Notification.Looping.Alarm9"
	LoopingAlarm10 Audio = "ms-winsoundevent:Notification.Looping.Alarm10"
	LoopingCall    Audio = "ms-winsoundevent:Notification.Looping.Call"
	LoopingCall2   Audio = "ms-winsoundevent:Notification.Looping.Call2"
	LoopingCall3   Audio = "ms-winsoundevent:Notification.Looping.Call3"
	LoopingCall4   Audio = "ms-winsoundevent:Notification.Looping.Call4"
	LoopingCall5   Audio = "ms-winsoundevent:Notification.Looping.Call5"
	LoopingCall6   Audio = "ms-winsoundevent:Notification.Looping.Call6"
	LoopingCall7   Audio = "ms-winsoundevent:Notification.Looping.Call7"
	LoopingCall8   Audio = "ms-winsoundevent:Notification.Looping.Call8"
	LoopingCall9   Audio = "ms-winsoundevent:Notification.Looping.Call9"
	LoopingCall10  Audio = "ms-winsoundevent:Notification.Looping.Call10"
	Silent         Audio = "silent"
)

const (
	Short Duration = "short"
	Long  Duration = "long"
)

const (
	Foreground ActivationType = "foreground"
	Background ActivationType = "background"
	Protocol   ActivationType = "protocol"
	System     ActivationType = "system"
)

const (
	DefaultScenario  Scenario = "default"
	Alarm            Scenario = "alarm"
	ReminderScenario Scenario = "reminder"
	IncomingCall     Scenario = "incomingCall"
)

func init() {
	toastTemplate = template.Must(
		template.New("toast").Parse(`
[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null
[Windows.UI.Notifications.ToastNotification, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null
[Windows.Data.Xml.Dom.XmlDocument, Windows.Data.Xml.Dom.XmlDocument, ContentType = WindowsRuntime] | Out-Null

$APP_ID = '{{if .AppID}}{{.AppID}}{{else}}Windows App{{end}}'

$template = @"
<toast activationType="{{.ActivationType}}" launch="{{.ActivationArguments}}" duration="{{.Duration}}">
    <visual>
        <binding template="ToastGeneric">
            {{if .Icon}}<image placement="appLogoOverride" {{if .HintCropCircle}}hint-crop="circle"{{end}} src="{{.Icon}}" />{{end}}
            {{if .Hero}}<image placement="hero" src="{{.Hero}}" />{{end}}
            {{if .InlineImage}}<image src="{{.InlineImage}}" />{{end}}
            {{if .Title}}<text><![CDATA[{{.Title}}]]></text>{{end}}
            {{if .Message}}<text><![CDATA[{{.Message}}]]></text>{{end}}
            {{if .Attribution}}<text><![CDATA[{{.Attribution}}]]></text>{{end}}
        </binding>
    </visual>
    {{if ne .Audio "silent"}}<audio src="{{.Audio}}" loop="{{.Loop}}" />{{else}}<audio silent="true" />{{end}}
    {{if .Actions}}
    <actions>
        {{range .Actions}}
        <action activationType="{{.ActivationType}}" content="{{.Content}}" arguments="{{.Arguments}}" />
        {{end}}
    </actions>
    {{end}}
</toast>
"@
$xml = New-Object Windows.Data.Xml.Dom.XmlDocument
$xml.LoadXml($template)
$toast = New-Object Windows.UI.Notifications.ToastNotification $xml
[Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier($APP_ID).Show($toast)
`))
}

type Action struct {
	ActivationType ActivationType
	Content        string
	Arguments      string
	ImageURI       string
}

type Notification struct {
	AppID               string
	Title               string
	Message             string
	Icon                string
	Hero                string
	InlineImage         string
	ActivationType      ActivationType
	Scenario            Scenario
	ActivationArguments string
	Actions             []Action
	Audio               Audio
	Loop                bool
	Duration            Duration
	HintCropCircle      bool
	Attribution         string
}

func (n *Notification) applyDefaults() {
	if n.ActivationType == "" {
		n.ActivationType = Protocol
	}
	if n.Scenario == "" {
		n.Scenario = DefaultScenario
	}
	if n.Audio == "" {
		n.Audio = Silent
	}
	if n.Duration == "" {
		n.Duration = Short
	}
}

func (n *Notification) buildXML() (*bytes.Buffer, error) {
	n.applyDefaults()

	s := new(bytes.Buffer)
	err := toastTemplate.Execute(s, n)
	if err != nil {
		return nil, err
	}
	return s, err
}

func (n *Notification) Send() error {
	xmlContent, err := n.buildXML()
	if err != nil {
		return err
	}

	cmd := exec.Command("PowerShell", "-Command", xmlContent.String())
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
