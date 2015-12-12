/* gonuts10seqB.go - https://groups.google.com/forum/?fromgroups#!topic/golang-nuts/tf4aDQ1Hn_c

NOTE: use NewMapXmlSeq() and mv.XmlSeqIndent() to preserve structure.

See data value at EOF - from: https://gist.github.com/suntong/e4dcdc6c85dcf769eec4

Objective:  into Comment.CommentText attribute value into Request.ReportingName attribute that immediately follows Comment.
*/

package main

import (
	"fmt"
	"github.com/clbanning/mxj"
)

func main() {
	// fmt.Println(string(data))
	m, err := mxj.NewMapXmlSeq(data)
	if err != nil {
		fmt.Println("NewMapXmlSeq err:", err)
		return
	}
	// fmt.Println(m.StringIndent())

	vals, err := m.ValuesForPath("WebTest.Items") //
	if err != nil {
		fmt.Println("ValuesForPath err:", err)
		return
	} else if len(vals) == 0 {
		fmt.Println("no WebTest.Items vals")
		return
	}
	var cmt, req []interface{}
	for _, v := range vals {
		vm, ok := v.(map[string]interface{})
		if !ok {
			fmt.Println("assertion failed")
			return
		}
		// get the Comment list
		cmt, ok = vm["Comment"].([]interface{})
		if !ok {
			continue
		}
		// get the Request list
		req, ok = vm["Request"].([]interface{})
		if !ok {
			continue
		}
		if cmt == nil || req == nil {
			break
		}

		// fmt.Println("Comment:", cmt)
		// fmt.Println("Request:", req)

		// Comment elements with #seq==n are followed by Request element with #seq==n+1.
		// For each Comment(n) extract the CommentText attribute value and use it to
		// set the ReportingName attribute value in Request(n+1).
		for _, v := range cmt {
			seq := v.(map[string]interface{})["#seq"].(int) // type is int
			var acmt string
			var ok bool
			for _, vv := range v.(map[string]interface{})["#attr"].([]interface{}) {
				if acmt, ok = vv.(map[string]interface{})["CommentText"].(string); ok {
					break
				}
			}
			if acmt == "" {
				fmt.Println("no CommentText value in Comment attributes")
			}
			// fmt.Println(seq, acmt)
			// find the request with the #seq==seq+1 value
			var r map[string]interface{}
			for _, vv := range req {
				r = vv.(map[string]interface{})
				if r["#seq"].(int) == seq+1 {
					break
				}
			}
			// fmt.Println(r)
			// loop through attributes to fine the ReportingName entry
			var rn map[string]interface{}
			for _, vv := range r["#attr"].([]interface{}) {
				rn = vv.(map[string]interface{})
				for key, _ := range rn {
					if key == "ReportingName" {
						goto gotReqName
					}
				}
			}
		gotReqName:
			if rn == nil { // shouldn't happen, but be safe
				continue
			}
			// set it to acmt
			// if you just want first 10 chars: rn["ReportingName"] = acmt[:10]
			rn["ReportingName"] = acmt
			// fmt.Println(r)
		}
	}

	// now do the same thing for the WebTest.Items.Request.TransactionTimer.Items
	// But because of #seq values can be the same for Comments and Requests across
	// Items entries, we will process each Item individually.
	vals, err = m.ValuesForPath("WebTest.Items.TransactionTimer.Items") //
	if err != nil {
		fmt.Println("ValuesForPath err:", err)
		return
	} else if len(vals) == 0 {
		fmt.Println("no WebTest.Items.TransactionTime.Items vals")
		return
	}
	for _, v := range vals {
		vm, ok := v.(map[string]interface{})
		if !ok {
			fmt.Println("assertion failed")
			return
		}
		fmt.Println("vm:", vm)

		// get the Comment list
		cmt, ok = vm["Comment"].([]interface{})
		if !ok {
			continue
		}
		// get the Request list
		req, ok = vm["Request"].([]interface{})
		if !ok {
			continue
		}
		if cmt == nil || req == nil {
			break
		}

		// fmt.Println("Comment:", cmt)
		// fmt.Println("Request:", req)

		// Comment elements with #seq==n are followed by Request element with #seq==n+1.
		// For each Comment(n) extract the CommentText attribute value and use it to
		// set the ReportingName attribute value in Request(n+1).
		for _, v := range cmt {
			seq := v.(map[string]interface{})["#seq"].(int) // type is int
			var acmt string
			var ok bool
			for _, vv := range v.(map[string]interface{})["#attr"].([]interface{}) {
				if acmt, ok = vv.(map[string]interface{})["CommentText"].(string); ok {
					break
				}
			}
			if acmt == "" {
				fmt.Println("no CommentText value in Comment attributes")
			}
			// fmt.Println(seq, acmt)
			// find the request with the #seq==seq+1 value
			var r map[string]interface{}
			for _, vv := range req {
				r = vv.(map[string]interface{})
				if r["#seq"].(int) == seq+1 {
					break
				}
			}
			// fmt.Println(r)
			// loop through attributes to fine the ReportingName entry
			var rn map[string]interface{}
			for _, vv := range r["#attr"].([]interface{}) {
				rn = vv.(map[string]interface{})
				for key, _ := range rn {
					if key == "ReportingName" {
						goto gotTReqName
					}
				}
			}
		gotTReqName:
			if rn == nil { // shouldn't happen, but be safe
				continue
			}
			// set it to acmt
			// if you just want first 10 chars: rn["ReportingName"] = acmt[:10]
			rn["ReportingName"] = acmt
			// fmt.Println(r)
		}
	}

	b, err := m.XmlSeqIndent("", "  ")
	if err != nil {
		fmt.Println("XmlIndent err:", err)
		return
	}
	fmt.Println(string(b))
}

var data = []byte(`
<?xml version="1.0" encoding="utf-8"?>
<WebTest Name="FirstAnonymousVisit" Id="ac766d08-f940-4b0a-b8f8-80675978894e" Owner="" Priority="0" Enabled="True" CssProjectStructure="" CssIteration="" Timeout="0" WorkItemIds="" xmlns="http://microsoft.com/schemas/VisualStudio/TeamTest/2010" Description="" CredentialUserName="" CredentialPassword="" PreAuthenticate="True" Proxy="" StopOnError="False" RecordedResultFile="">
  <Items>
    <Comment CommentText="Visit Homepage and ensure new page setup is created" />
    <Request Method="GET" Version="1.1" Url="{{Config.TestParameters.ServerURL}}/Default.aspx" ThinkTime="0" Timeout="300" ParseDependentRequests="False" FollowRedirects="True" RecordResult="True" Cache="False" ResponseTimeGoal="0.5" Encoding="utf-8" ExpectedHttpStatusCode="0" ExpectedResponseUrl="" ReportingName="">
      <ValidationRules>
        <ValidationRule Classname="Dropthings.Test.Rules.CookieValidationRule, Dropthings.Test, Version=1.0.0.0, Culture=neutral, PublicKeyToken=null" DisplayName="Check Cookie From Response" Description="" Level="High" ExectuionOrder="BeforeDependents">
          <RuleParameters>
            <RuleParameter Name="StopOnError" Value="False" />
            <RuleParameter Name="CookieValueToMatch" Value="" />
            <RuleParameter Name="MatchValue" Value="False" />
            <RuleParameter Name="Exists" Value="True" />
            <RuleParameter Name="CookieName" Value="{{Config.TestParameters.AnonCookieName}}" />
            <RuleParameter Name="IsPersistent" Value="True" />
            <RuleParameter Name="Domain" Value="" />
            <RuleParameter Name="Index" Value="0" />
          </RuleParameters>
        </ValidationRule>
        <ValidationRule Classname="Dropthings.Test.Rules.CookieValidationRule, Dropthings.Test, Version=1.0.0.0, Culture=neutral, PublicKeyToken=null" DisplayName="Check Cookie From Response" Description="" Level="High" ExectuionOrder="BeforeDependents">
          <RuleParameters>
            <RuleParameter Name="StopOnError" Value="False" />
            <RuleParameter Name="CookieValueToMatch" Value="" />
            <RuleParameter Name="MatchValue" Value="False" />
            <RuleParameter Name="Exists" Value="False" />
            <RuleParameter Name="CookieName" Value="{{Config.TestParameters.SessionCookieName}}" />
            <RuleParameter Name="IsPersistent" Value="False" />
            <RuleParameter Name="Domain" Value="" />
            <RuleParameter Name="Index" Value="0" />
          </RuleParameters>
        </ValidationRule>
        <ValidationRule Classname="Dropthings.Test.Rules.CacheHeaderValidation, Dropthings.Test, Version=1.0.0.0, Culture=neutral, PublicKeyToken=null" DisplayName="Cache Header Validation" Description="" Level="High" ExectuionOrder="BeforeDependents">
          <RuleParameters>
            <RuleParameter Name="Enabled" Value="True" />
            <RuleParameter Name="DifferenceThresholdSec" Value="0" />
            <RuleParameter Name="CacheControlPrivate" Value="False" />
            <RuleParameter Name="CacheControlPublic" Value="False" />
            <RuleParameter Name="CacheControlNoCache" Value="True" />
            <RuleParameter Name="ExpiresAfterSeconds" Value="0" />
            <RuleParameter Name="StopOnError" Value="False" />
          </RuleParameters>
        </ValidationRule>
        <ValidationRule Classname="Microsoft.VisualStudio.TestTools.WebTesting.Rules.ValidationRuleFindText, Microsoft.VisualStudio.QualityTools.WebTestFramework, Version=10.0.0.0, Culture=neutral, PublicKeyToken=b03f5f7f11d50a3a" DisplayName="Find Text" Description="Verifies the existence of the specified text in the response." Level="High" ExectuionOrder="BeforeDependents">
          <RuleParameters>
            <RuleParameter Name="FindText" Value="How to of the Day" />
            <RuleParameter Name="IgnoreCase" Value="False" />
            <RuleParameter Name="UseRegularExpression" Value="False" />
            <RuleParameter Name="PassIfTextFound" Value="True" />
          </RuleParameters>
        </ValidationRule>
        <ValidationRule Classname="Microsoft.VisualStudio.TestTools.WebTesting.Rules.ValidationRuleFindText, Microsoft.VisualStudio.QualityTools.WebTestFramework, Version=10.0.0.0, Culture=neutral, PublicKeyToken=b03f5f7f11d50a3a" DisplayName="Find Text" Description="Verifies the existence of the specified text in the response." Level="High" ExectuionOrder="BeforeDependents">
          <RuleParameters>
            <RuleParameter Name="FindText" Value="Weather" />
            <RuleParameter Name="IgnoreCase" Value="False" />
            <RuleParameter Name="UseRegularExpression" Value="False" />
            <RuleParameter Name="PassIfTextFound" Value="True" />
          </RuleParameters>
        </ValidationRule>
        <ValidationRule Classname="Microsoft.VisualStudio.TestTools.WebTesting.Rules.ValidationRuleFindText, Microsoft.VisualStudio.QualityTools.WebTestFramework, Version=10.0.0.0, Culture=neutral, PublicKeyToken=b03f5f7f11d50a3a" DisplayName="Find Text" Description="Verifies the existence of the specified text in the response." Level="High" ExectuionOrder="BeforeDependents">
          <RuleParameters>
            <RuleParameter Name="FindText" Value="All rights reserved" />
            <RuleParameter Name="IgnoreCase" Value="False" />
            <RuleParameter Name="UseRegularExpression" Value="False" />
            <RuleParameter Name="PassIfTextFound" Value="True" />
          </RuleParameters>
        </ValidationRule>
      </ValidationRules>
    </Request>
    <TransactionTimer Name="Show Hide Widget List">
      <Items>
        <Comment CommentText="Show Widget List and expect Widget List to produce BBC Word widget link" />
        <Request Method="GET" Version="1.1" Url="{{Config.TestParameters.ServerURL}}/Default.aspx" ThinkTime="0" Timeout="300" ParseDependentRequests="False" FollowRedirects="True" RecordResult="True" Cache="False" ResponseTimeGoal="0.5" Encoding="utf-8" ExpectedHttpStatusCode="0" ExpectedResponseUrl="" ReportingName="">
          <ValidationRules>
            <ValidationRule Classname="Microsoft.VisualStudio.TestTools.WebTesting.Rules.ValidationRuleFindText, Microsoft.VisualStudio.QualityTools.WebTestFramework, Version=10.0.0.0, Culture=neutral, PublicKeyToken=b03f5f7f11d50a3a" DisplayName="Find Text" Description="Verifies the existence of the specified text in the response." Level="High" ExectuionOrder="BeforeDependents">
              <RuleParameters>
                <RuleParameter Name="FindText" Value="BBC World" />
                <RuleParameter Name="IgnoreCase" Value="False" />
                <RuleParameter Name="UseRegularExpression" Value="False" />
                <RuleParameter Name="PassIfTextFound" Value="True" />
              </RuleParameters>
            </ValidationRule>
          </ValidationRules>
          <RequestPlugins>
            <RequestPlugin Classname="Dropthings.Test.Plugin.AsyncPostbackRequestPlugin, Dropthings.Test, Version=1.0.0.0, Culture=neutral, PublicKeyToken=null" DisplayName="AsyncPostbackRequestPlugin" Description="">
              <RuleParameters>
                <RuleParameter Name="ControlName" Value="TabControlPanel$ShowAddContentPanel" />
                <RuleParameter Name="UpdatePanelName" Value="{{$UPDATEPANEL.OnPageMenuUpdatePanel.1}}" />
              </RuleParameters>
            </RequestPlugin>
          </RequestPlugins>
        </Request>
        <Comment CommentText="Hide Widget List and expect the outpu does not have the BBC World Widget" />
        <Request Method="GET" Version="1.1" Url="{{Config.TestParameters.ServerURL}}/Default.aspx" ThinkTime="0" Timeout="300" ParseDependentRequests="False" FollowRedirects="True" RecordResult="True" Cache="False" ResponseTimeGoal="0.5" Encoding="utf-8" ExpectedHttpStatusCode="0" ExpectedResponseUrl="" ReportingName="">
          <ValidationRules>
            <ValidationRule Classname="Microsoft.VisualStudio.TestTools.WebTesting.Rules.ValidationRuleFindText, Microsoft.VisualStudio.QualityTools.WebTestFramework, Version=10.0.0.0, Culture=neutral, PublicKeyToken=b03f5f7f11d50a3a" DisplayName="Find Text" Description="Verifies the existence of the specified text in the response." Level="High" ExectuionOrder="BeforeDependents">
              <RuleParameters>
                <RuleParameter Name="FindText" Value="TabControlPanel$ShowAddContentPanel" />
                <RuleParameter Name="IgnoreCase" Value="False" />
                <RuleParameter Name="UseRegularExpression" Value="False" />
                <RuleParameter Name="PassIfTextFound" Value="True" />
              </RuleParameters>
            </ValidationRule>
          </ValidationRules>
          <RequestPlugins>
            <RequestPlugin Classname="Dropthings.Test.Plugin.AsyncPostbackRequestPlugin, Dropthings.Test, Version=1.0.0.0, Culture=neutral, PublicKeyToken=null" DisplayName="AsyncPostbackRequestPlugin" Description="">
              <RuleParameters>
                <RuleParameter Name="ControlName" Value="TabControlPanel$HideAddContentPanel" />
                <RuleParameter Name="UpdatePanelName" Value="{{$UPDATEPANEL.OnPageMenuUpdatePanel.1}}" />
              </RuleParameters>
            </RequestPlugin>
          </RequestPlugins>
        </Request>
      </Items>
    </TransactionTimer>
    <Request Method="GET" Version="1.1" Url="{{Config.TestParameters.ServerURL}}/API/Proxy.svc/ajax/GetRss?url=%22http%3A%2F%2Ffeeds.feedburner.com%2FOmarAlZabirBlog%22&amp;count=10&amp;cacheDuration=10" ThinkTime="0" Timeout="300" ParseDependentRequests="True" FollowRedirects="True" RecordResult="True" Cache="False" ResponseTimeGoal="0" Encoding="utf-8" ExpectedHttpStatusCode="0" ExpectedResponseUrl="" ReportingName="">
      <ValidationRules>
        <ValidationRule Classname="Microsoft.VisualStudio.TestTools.WebTesting.Rules.ValidationRuleFindText, Microsoft.VisualStudio.QualityTools.WebTestFramework, Version=10.0.0.0, Culture=neutral, PublicKeyToken=b03f5f7f11d50a3a" DisplayName="Find Text" Description="Verifies the existence of the specified text in the response." Level="High" ExectuionOrder="BeforeDependents">
          <RuleParameters>
            <RuleParameter Name="FindText" Value="{&quot;d&quot;:[{&quot;__type&quot;:&quot;RssItem:#Dropthings.Web.Util&quot;" />
            <RuleParameter Name="IgnoreCase" Value="False" />
            <RuleParameter Name="UseRegularExpression" Value="False" />
            <RuleParameter Name="PassIfTextFound" Value="True" />
          </RuleParameters>
        </ValidationRule>
      </ValidationRules>
    </Request>
    <Request Method="GET" Version="1.1" Url="{{Config.TestParameters.ServerURL}}/API/Proxy.svc/ajax/GetUrl?url=%22http%3A%2F%2Ffeeds.feedburner.com%2FOmarAlZabirBlog%22&amp;cacheDuration=10" ThinkTime="0" Timeout="300" ParseDependentRequests="True" FollowRedirects="True" RecordResult="True" Cache="False" ResponseTimeGoal="0" Encoding="utf-8" ExpectedHttpStatusCode="0" ExpectedResponseUrl="" ReportingName="">
      <ValidationRules>
        <ValidationRule Classname="Microsoft.VisualStudio.TestTools.WebTesting.Rules.ValidationRuleFindText, Microsoft.VisualStudio.QualityTools.WebTestFramework, Version=10.0.0.0, Culture=neutral, PublicKeyToken=b03f5f7f11d50a3a" DisplayName="Find Text" Description="Verifies the existence of the specified text in the response." Level="High" ExectuionOrder="BeforeDependents">
          <RuleParameters>
            <RuleParameter Name="FindText" Value="&lt;channel&gt;" />
            <RuleParameter Name="IgnoreCase" Value="False" />
            <RuleParameter Name="UseRegularExpression" Value="False" />
            <RuleParameter Name="PassIfTextFound" Value="True" />
          </RuleParameters>
        </ValidationRule>
      </ValidationRules>
    </Request>
    <TransactionTimer Name="Edit Collapse Expand Widget">
      <Items>
        <Comment CommentText="Click edit on first widget &quot;How to of the Day&quot; and expect URL textbox to be present with Feed Url" />
        <Request Method="GET" Version="1.1" Url="{{Config.TestParameters.ServerURL}}/Default.aspx" ThinkTime="0" Timeout="300" ParseDependentRequests="False" FollowRedirects="True" RecordResult="True" Cache="False" ResponseTimeGoal="0.5" Encoding="utf-8" ExpectedHttpStatusCode="0" ExpectedResponseUrl="" ReportingName="">
          <ValidationRules>
            <ValidationRule Classname="Microsoft.VisualStudio.TestTools.WebTesting.Rules.ValidationRuleRequiredAttributeValue, Microsoft.VisualStudio.QualityTools.WebTestFramework, Version=10.0.0.0, Culture=neutral, PublicKeyToken=b03f5f7f11d50a3a" DisplayName="Required Attribute Value" Description="Verifies the existence of a specified HTML tag that contains an attribute with a specified value." Level="High" ExectuionOrder="BeforeDependents">
              <RuleParameters>
                <RuleParameter Name="TagName" Value="input" />
                <RuleParameter Name="AttributeName" Value="value" />
                <RuleParameter Name="MatchAttributeName" Value="" />
                <RuleParameter Name="MatchAttributeValue" Value="" />
                <RuleParameter Name="ExpectedValue" Value="http://www.wikihow.com/feed.rss" />
                <RuleParameter Name="IgnoreCase" Value="False" />
                <RuleParameter Name="Index" Value="-1" />
              </RuleParameters>
            </ValidationRule>
          </ValidationRules>
          <ExtractionRules>
            <ExtractionRule Classname="Dropthings.Test.Rules.ExtractFormElements, Dropthings.Test, Version=1.0.0.0, Culture=neutral, PublicKeyToken=null" VariableName="" DisplayName="Extract Form Elements" Description="">
              <RuleParameters>
                <RuleParameter Name="ContextParameterName" Value="" />
              </RuleParameters>
            </ExtractionRule>
          </ExtractionRules>
          <RequestPlugins>
            <RequestPlugin Classname="Dropthings.Test.Plugin.AsyncPostbackRequestPlugin, Dropthings.Test, Version=1.0.0.0, Culture=neutral, PublicKeyToken=null" DisplayName="AsyncPostbackRequestPlugin" Description="">
              <RuleParameters>
                <RuleParameter Name="ControlName" Value="{{$POSTBACK.EditWidget.1}}" />
                <RuleParameter Name="UpdatePanelName" Value="{{$UPDATEPANEL.WidgetHeaderUpdatePanel.1}}" />
              </RuleParameters>
            </RequestPlugin>
          </RequestPlugins>
        </Request>
        <Comment CommentText="Change the Feed Count Dropdown list to 10 and expect 10 Feed Link controls are generated" />
        <Request Method="POST" Version="1.1" Url="{{Config.TestParameters.ServerURL}}/Default.aspx" ThinkTime="0" Timeout="300" ParseDependentRequests="False" FollowRedirects="True" RecordResult="True" Cache="False" ResponseTimeGoal="0.5" Encoding="utf-8" ExpectedHttpStatusCode="0" ExpectedResponseUrl="" ReportingName="">
          <ValidationRules>
            <ValidationRule Classname="Microsoft.VisualStudio.TestTools.WebTesting.Rules.ValidationRuleFindText, Microsoft.VisualStudio.QualityTools.WebTestFramework, Version=10.0.0.0, Culture=neutral, PublicKeyToken=b03f5f7f11d50a3a" DisplayName="Find Text" Description="Verifies the existence of the specified text in the response." Level="High" ExectuionOrder="BeforeDependents">
              <RuleParameters>
                <RuleParameter Name="FindText" Value="FeedList_ctl09_FeedLink" />
                <RuleParameter Name="IgnoreCase" Value="False" />
                <RuleParameter Name="UseRegularExpression" Value="False" />
                <RuleParameter Name="PassIfTextFound" Value="True" />
              </RuleParameters>
            </ValidationRule>
          </ValidationRules>
          <RequestPlugins>
            <RequestPlugin Classname="Dropthings.Test.Plugin.AsyncPostbackRequestPlugin, Dropthings.Test, Version=1.0.0.0, Culture=neutral, PublicKeyToken=null" DisplayName="AsyncPostbackRequestPlugin" Description="">
              <RuleParameters>
                <RuleParameter Name="ControlName" Value="{{$POSTBACK.CancelEditWidget.1}}" />
                <RuleParameter Name="UpdatePanelName" Value="{{$UPDATEPANEL.WidgetHeaderUpdatePanel.1}}" />
              </RuleParameters>
            </RequestPlugin>
          </RequestPlugins>
          <FormPostHttpBody>
            <FormPostParameter Name="{{$INPUT.FeedUrl.1}}" Value="http://www.wikihow.com/feed.rss" RecordedValue="" CorrelationBinding="" UrlEncode="True" />
            <FormPostParameter Name="{{$SELECT.FeedCountDropDownList.1}}" Value="10" RecordedValue="" CorrelationBinding="" UrlEncode="True" />
          </FormPostHttpBody>
        </Request>
        <Comment CommentText="Delete the How to of the Day widget and expect it's not found from response" />
        <Request Method="GET" Version="1.1" Url="{{Config.TestParameters.ServerURL}}/Default.aspx" ThinkTime="0" Timeout="300" ParseDependentRequests="False" FollowRedirects="True" RecordResult="True" Cache="False" ResponseTimeGoal="0.5" Encoding="utf-8" ExpectedHttpStatusCode="0" ExpectedResponseUrl="" ReportingName="">
          <ValidationRules>
            <ValidationRule Classname="Microsoft.VisualStudio.TestTools.WebTesting.Rules.ValidationRuleFindText, Microsoft.VisualStudio.QualityTools.WebTestFramework, Version=10.0.0.0, Culture=neutral, PublicKeyToken=b03f5f7f11d50a3a" DisplayName="Find Text" Description="Verifies the existence of the specified text in the response." Level="High" ExectuionOrder="BeforeDependents">
              <RuleParameters>
                <RuleParameter Name="FindText" Value="How to of the Day" />
                <RuleParameter Name="IgnoreCase" Value="False" />
                <RuleParameter Name="UseRegularExpression" Value="False" />
                <RuleParameter Name="PassIfTextFound" Value="False" />
              </RuleParameters>
            </ValidationRule>
          </ValidationRules>
          <RequestPlugins>
            <RequestPlugin Classname="Dropthings.Test.Plugin.AsyncPostbackRequestPlugin, Dropthings.Test, Version=1.0.0.0, Culture=neutral, PublicKeyToken=null" DisplayName="AsyncPostbackRequestPlugin" Description="">
              <RuleParameters>
                <RuleParameter Name="ControlName" Value="{{$POSTBACK.CloseWidget.1}}" />
                <RuleParameter Name="UpdatePanelName" Value="{{$UPDATEPANEL.WidgetHeaderUpdatePanel.1}}" />
              </RuleParameters>
            </RequestPlugin>
          </RequestPlugins>
        </Request>
      </Items>
    </TransactionTimer>
    <TransactionTimer Name="Add New Widget">
      <Items>
        <Comment CommentText="Show widget list and expect Digg to be there" />
        <Request Method="GET" Version="1.1" Url="{{Config.TestParameters.ServerURL}}/Default.aspx" ThinkTime="0" Timeout="300" ParseDependentRequests="False" FollowRedirects="True" RecordResult="True" Cache="False" ResponseTimeGoal="0.5" Encoding="utf-8" ExpectedHttpStatusCode="0" ExpectedResponseUrl="" ReportingName="">
          <ValidationRules>
            <ValidationRule Classname="Microsoft.VisualStudio.TestTools.WebTesting.Rules.ValidationRuleFindText, Microsoft.VisualStudio.QualityTools.WebTestFramework, Version=10.0.0.0, Culture=neutral, PublicKeyToken=b03f5f7f11d50a3a" DisplayName="Find Text" Description="Verifies the existence of the specified text in the response." Level="High" ExectuionOrder="BeforeDependents">
              <RuleParameters>
                <RuleParameter Name="FindText" Value="Digg" />
                <RuleParameter Name="IgnoreCase" Value="False" />
                <RuleParameter Name="UseRegularExpression" Value="False" />
                <RuleParameter Name="PassIfTextFound" Value="True" />
              </RuleParameters>
            </ValidationRule>
          </ValidationRules>
          <RequestPlugins>
            <RequestPlugin Classname="Dropthings.Test.Plugin.AsyncPostbackRequestPlugin, Dropthings.Test, Version=1.0.0.0, Culture=neutral, PublicKeyToken=null" DisplayName="AsyncPostbackRequestPlugin" Description="">
              <RuleParameters>
                <RuleParameter Name="ControlName" Value="TabControlPanel$ShowAddContentPanel" />
                <RuleParameter Name="UpdatePanelName" Value="{{$UPDATEPANEL.OnPageMenuUpdatePanel.1}}" />
              </RuleParameters>
            </RequestPlugin>
          </RequestPlugins>
        </Request>
        <Comment CommentText="Add New Widget" />
        <Request Method="GET" Version="1.1" Url="{{Config.TestParameters.ServerURL}}/Default.aspx" ThinkTime="0" Timeout="300" ParseDependentRequests="False" FollowRedirects="True" RecordResult="True" Cache="False" ResponseTimeGoal="0.5" Encoding="utf-8" ExpectedHttpStatusCode="0" ExpectedResponseUrl="" ReportingName="">
          <ValidationRules>
            <ValidationRule Classname="Microsoft.VisualStudio.TestTools.WebTesting.Rules.ValidationRuleFindText, Microsoft.VisualStudio.QualityTools.WebTestFramework, Version=10.0.0.0, Culture=neutral, PublicKeyToken=b03f5f7f11d50a3a" DisplayName="Find Text" Description="Verifies the existence of the specified text in the response." Level="High" ExectuionOrder="BeforeDependents">
              <RuleParameters>
                <RuleParameter Name="FindText" Value="Digg" />
                <RuleParameter Name="IgnoreCase" Value="False" />
                <RuleParameter Name="UseRegularExpression" Value="False" />
                <RuleParameter Name="PassIfTextFound" Value="True" />
              </RuleParameters>
            </ValidationRule>
          </ValidationRules>
          <RequestPlugins>
            <RequestPlugin Classname="Dropthings.Test.Plugin.AsyncPostbackRequestPlugin, Dropthings.Test, Version=1.0.0.0, Culture=neutral, PublicKeyToken=null" DisplayName="AsyncPostbackRequestPlugin" Description="">
              <RuleParameters>
                <RuleParameter Name="ControlName" Value="{{$POSTBACK.AddWidget.1}}" />
                <RuleParameter Name="UpdatePanelName" Value="{{$UPDATEPANEL.OnPageMenuUpdatePanel.1}}" />
              </RuleParameters>
            </RequestPlugin>
          </RequestPlugins>
        </Request>
        <Comment CommentText="Delete the newly added widget" />
        <Request Method="GET" Version="1.1" Url="{{Config.TestParameters.ServerURL}}/Default.aspx" ThinkTime="0" Timeout="300" ParseDependentRequests="False" FollowRedirects="True" RecordResult="True" Cache="False" ResponseTimeGoal="0.5" Encoding="utf-8" ExpectedHttpStatusCode="0" ExpectedResponseUrl="" ReportingName="">
          <ValidationRules>
            <ValidationRule Classname="Microsoft.VisualStudio.TestTools.WebTesting.Rules.ValidationRuleFindText, Microsoft.VisualStudio.QualityTools.WebTestFramework, Version=10.0.0.0, Culture=neutral, PublicKeyToken=b03f5f7f11d50a3a" DisplayName="Find Text" Description="Verifies the existence of the specified text in the response." Level="High" ExectuionOrder="BeforeDependents">
              <RuleParameters>
                <RuleParameter Name="FindText" Value="How to of the Day" />
                <RuleParameter Name="IgnoreCase" Value="False" />
                <RuleParameter Name="UseRegularExpression" Value="False" />
                <RuleParameter Name="PassIfTextFound" Value="False" />
              </RuleParameters>
            </ValidationRule>
          </ValidationRules>
          <RequestPlugins>
            <RequestPlugin Classname="Dropthings.Test.Plugin.AsyncPostbackRequestPlugin, Dropthings.Test, Version=1.0.0.0, Culture=neutral, PublicKeyToken=null" DisplayName="AsyncPostbackRequestPlugin" Description="">
              <RuleParameters>
                <RuleParameter Name="ControlName" Value="{{$POSTBACK.CloseWidget.1}}" />
                <RuleParameter Name="UpdatePanelName" Value="{{$UPDATEPANEL.WidgetHeaderUpdatePanel.1}}" />
              </RuleParameters>
            </RequestPlugin>
          </RequestPlugins>
        </Request>
      </Items>
    </TransactionTimer>
    <Comment CommentText="Revisit and ensure the Digg widget exists and How to widget does not exist" />
    <Request Method="GET" Version="1.1" Url="{{Config.TestParameters.ServerURL}}/Default.aspx" ThinkTime="0" Timeout="300" ParseDependentRequests="False" FollowRedirects="True" RecordResult="True" Cache="False" ResponseTimeGoal="0.5" Encoding="utf-8" ExpectedHttpStatusCode="0" ExpectedResponseUrl="" ReportingName="">
      <ValidationRules>
        <ValidationRule Classname="Dropthings.Test.Rules.CookieValidationRule, Dropthings.Test, Version=1.0.0.0, Culture=neutral, PublicKeyToken=null" DisplayName="Check Cookie From Response" Description="" Level="High" ExectuionOrder="BeforeDependents">
          <RuleParameters>
            <RuleParameter Name="StopOnError" Value="False" />
            <RuleParameter Name="CookieValueToMatch" Value="" />
            <RuleParameter Name="MatchValue" Value="False" />
            <RuleParameter Name="Exists" Value="False" />
            <RuleParameter Name="CookieName" Value="{{Config.TestParameters.AnonCookieName}}" />
            <RuleParameter Name="IsPersistent" Value="True" />
            <RuleParameter Name="Domain" Value="" />
            <RuleParameter Name="Index" Value="0" />
          </RuleParameters>
        </ValidationRule>
        <ValidationRule Classname="Dropthings.Test.Rules.CookieValidationRule, Dropthings.Test, Version=1.0.0.0, Culture=neutral, PublicKeyToken=null" DisplayName="Check Cookie From Response" Description="" Level="High" ExectuionOrder="BeforeDependents">
          <RuleParameters>
            <RuleParameter Name="StopOnError" Value="False" />
            <RuleParameter Name="CookieValueToMatch" Value="" />
            <RuleParameter Name="MatchValue" Value="False" />
            <RuleParameter Name="Exists" Value="False" />
            <RuleParameter Name="CookieName" Value="{{Config.TestParameters.SessionCookieName}}" />
            <RuleParameter Name="IsPersistent" Value="False" />
            <RuleParameter Name="Domain" Value="" />
            <RuleParameter Name="Index" Value="0" />
          </RuleParameters>
        </ValidationRule>
        <ValidationRule Classname="Dropthings.Test.Rules.CacheHeaderValidation, Dropthings.Test, Version=1.0.0.0, Culture=neutral, PublicKeyToken=null" DisplayName="Cache Header Validation" Description="" Level="High" ExectuionOrder="BeforeDependents">
          <RuleParameters>
            <RuleParameter Name="Enabled" Value="True" />
            <RuleParameter Name="DifferenceThresholdSec" Value="0" />
            <RuleParameter Name="CacheControlPrivate" Value="False" />
            <RuleParameter Name="CacheControlPublic" Value="False" />
            <RuleParameter Name="CacheControlNoCache" Value="True" />
            <RuleParameter Name="ExpiresAfterSeconds" Value="0" />
            <RuleParameter Name="StopOnError" Value="False" />
          </RuleParameters>
        </ValidationRule>
        <ValidationRule Classname="Microsoft.VisualStudio.TestTools.WebTesting.Rules.ValidationRuleFindText, Microsoft.VisualStudio.QualityTools.WebTestFramework, Version=10.0.0.0, Culture=neutral, PublicKeyToken=b03f5f7f11d50a3a" DisplayName="Find Text" Description="Verifies the existence of the specified text in the response." Level="High" ExectuionOrder="BeforeDependents">
          <RuleParameters>
            <RuleParameter Name="FindText" Value="How to of the Day" />
            <RuleParameter Name="IgnoreCase" Value="False" />
            <RuleParameter Name="UseRegularExpression" Value="False" />
            <RuleParameter Name="PassIfTextFound" Value="False" />
          </RuleParameters>
        </ValidationRule>
        <ValidationRule Classname="Microsoft.VisualStudio.TestTools.WebTesting.Rules.ValidationRuleFindText, Microsoft.VisualStudio.QualityTools.WebTestFramework, Version=10.0.0.0, Culture=neutral, PublicKeyToken=b03f5f7f11d50a3a" DisplayName="Find Text" Description="Verifies the existence of the specified text in the response." Level="High" ExectuionOrder="BeforeDependents">
          <RuleParameters>
            <RuleParameter Name="FindText" Value="Digg" />
            <RuleParameter Name="IgnoreCase" Value="False" />
            <RuleParameter Name="UseRegularExpression" Value="False" />
            <RuleParameter Name="PassIfTextFound" Value="True" />
          </RuleParameters>
        </ValidationRule>
        <ValidationRule Classname="Microsoft.VisualStudio.TestTools.WebTesting.Rules.ValidationRuleFindText, Microsoft.VisualStudio.QualityTools.WebTestFramework, Version=10.0.0.0, Culture=neutral, PublicKeyToken=b03f5f7f11d50a3a" DisplayName="Find Text" Description="Verifies the existence of the specified text in the response." Level="High" ExectuionOrder="BeforeDependents">
          <RuleParameters>
            <RuleParameter Name="FindText" Value="All rights reserved" />
            <RuleParameter Name="IgnoreCase" Value="False" />
            <RuleParameter Name="UseRegularExpression" Value="False" />
            <RuleParameter Name="PassIfTextFound" Value="True" />
          </RuleParameters>
        </ValidationRule>
      </ValidationRules>
    </Request>
    <Comment CommentText="- Logout and ensure Anon Cookie is set to expire" />
    <Request Method="GET" Version="1.1" Url="{{Config.TestParameters.ServerURL}}/Logout.ashx" ThinkTime="0" Timeout="300" ParseDependentRequests="False" FollowRedirects="False" RecordResult="True" Cache="False" ResponseTimeGoal="0.5" Encoding="utf-8" ExpectedHttpStatusCode="302" ExpectedResponseUrl="" ReportingName="">
      <ValidationRules>
        <ValidationRule Classname="Dropthings.Test.Rules.CookieSetToExpire, Dropthings.Test, Version=1.0.0.0, Culture=neutral, PublicKeyToken=null" DisplayName="Ensure Cookie Set to Expire" Description="" Level="High" ExectuionOrder="BeforeDependents">
          <RuleParameters>
            <RuleParameter Name="CookieName" Value="{{Config.TestParameters.AnonCookieName}}" />
            <RuleParameter Name="Domain" Value="" />
            <RuleParameter Name="StopOnError" Value="False" />
          </RuleParameters>
        </ValidationRule>
      </ValidationRules>
    </Request>
  </Items>
  <DataSources>
    <DataSource Name="Config" Provider="Microsoft.VisualStudio.TestTools.DataSource.XML" Connection="|DataDirectory|\Config\TestParameters.xml">
      <Tables>
        <DataSourceTable Name="TestParameters" SelectColumns="SelectOnlyBoundColumns" AccessMethod="Sequential" />
      </Tables>
    </DataSource>
  </DataSources>
  <ValidationRules>
    <ValidationRule Classname="Microsoft.VisualStudio.TestTools.WebTesting.Rules.ValidateResponseUrl, Microsoft.VisualStudio.QualityTools.WebTestFramework, Version=10.0.0.0, Culture=neutral, PublicKeyToken=b03f5f7f11d50a3a" DisplayName="Response URL" Description="Validates that the response URL after redirects are followed is the same as the recorded response URL.  QueryString parameters are ignored." Level="Low" ExectuionOrder="BeforeDependents" />
  </ValidationRules>
  <WebTestPlugins>
    <WebTestPlugin Classname="Dropthings.Test.Plugin.ASPNETWebTestPlugin, Dropthings.Test, Version=1.0.0.0, Culture=neutral, PublicKeyToken=null" DisplayName="ASPNETWebTestPlugin" Description="" />
  </WebTestPlugins>
</WebTest>
`)