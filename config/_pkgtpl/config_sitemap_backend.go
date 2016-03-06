// +build ignore

package sitemap

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/model"
)

// Backend will be initialized in the init() function together with ConfigStructure.
var Backend *PkgBackend

// PkgBackend just exported for the sake of documentation. See fields
// for more information. The PkgBackend handles the reading and writing
// of configuration values within this package.
type PkgBackend struct {
	model.PkgBackend
	// SitemapCategoryChangefreq => Frequency.
	// Path: sitemap/category/changefreq
	// SourceModel: Magento\Sitemap\Model\Config\Source\Frequency
	SitemapCategoryChangefreq model.Str

	// SitemapCategoryPriority => Priority.
	// Valid values range from 0.0 to 1.0.
	// Path: sitemap/category/priority
	// BackendModel: Magento\Sitemap\Model\Config\Backend\Priority
	SitemapCategoryPriority model.Str

	// SitemapProductChangefreq => Frequency.
	// Path: sitemap/product/changefreq
	// SourceModel: Magento\Sitemap\Model\Config\Source\Frequency
	SitemapProductChangefreq model.Str

	// SitemapProductPriority => Priority.
	// Valid values range from 0.0 to 1.0.
	// Path: sitemap/product/priority
	// BackendModel: Magento\Sitemap\Model\Config\Backend\Priority
	SitemapProductPriority model.Str

	// SitemapProductImageInclude => Add Images into Sitemap.
	// Path: sitemap/product/image_include
	// SourceModel: Magento\Sitemap\Model\Source\Product\Image\IncludeImage
	SitemapProductImageInclude model.Str

	// SitemapPageChangefreq => Frequency.
	// Path: sitemap/page/changefreq
	// SourceModel: Magento\Sitemap\Model\Config\Source\Frequency
	SitemapPageChangefreq model.Str

	// SitemapPagePriority => Priority.
	// Valid values range from 0.0 to 1.0.
	// Path: sitemap/page/priority
	// BackendModel: Magento\Sitemap\Model\Config\Backend\Priority
	SitemapPagePriority model.Str

	// SitemapGenerateEnabled => Enabled.
	// Path: sitemap/generate/enabled
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SitemapGenerateEnabled model.Bool

	// SitemapGenerateErrorEmail => Error Email Recipient.
	// Path: sitemap/generate/error_email
	SitemapGenerateErrorEmail model.Str

	// SitemapGenerateErrorEmailIdentity => Error Email Sender.
	// Path: sitemap/generate/error_email_identity
	// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
	SitemapGenerateErrorEmailIdentity model.Str

	// SitemapGenerateErrorEmailTemplate => Error Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: sitemap/generate/error_email_template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	SitemapGenerateErrorEmailTemplate model.Str

	// SitemapGenerateFrequency => Frequency.
	// Path: sitemap/generate/frequency
	// BackendModel: Magento\Cron\Model\Config\Backend\Sitemap
	// SourceModel: Magento\Cron\Model\Config\Source\Frequency
	SitemapGenerateFrequency model.Str

	// SitemapGenerateTime => Start Time.
	// Path: sitemap/generate/time
	SitemapGenerateTime model.Str

	// SitemapLimitMaxLines => Maximum No of URLs Per File.
	// Path: sitemap/limit/max_lines
	SitemapLimitMaxLines model.Str

	// SitemapLimitMaxFileSize => Maximum File Size.
	// File size in bytes.
	// Path: sitemap/limit/max_file_size
	SitemapLimitMaxFileSize model.Str

	// SitemapSearchEnginesSubmissionRobots => Enable Submission to Robots.txt.
	// Path: sitemap/search_engines/submission_robots
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SitemapSearchEnginesSubmissionRobots model.Bool
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.SitemapCategoryChangefreq = model.NewStr(`sitemap/category/changefreq`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SitemapCategoryPriority = model.NewStr(`sitemap/category/priority`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SitemapProductChangefreq = model.NewStr(`sitemap/product/changefreq`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SitemapProductPriority = model.NewStr(`sitemap/product/priority`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SitemapProductImageInclude = model.NewStr(`sitemap/product/image_include`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SitemapPageChangefreq = model.NewStr(`sitemap/page/changefreq`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SitemapPagePriority = model.NewStr(`sitemap/page/priority`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SitemapGenerateEnabled = model.NewBool(`sitemap/generate/enabled`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SitemapGenerateErrorEmail = model.NewStr(`sitemap/generate/error_email`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SitemapGenerateErrorEmailIdentity = model.NewStr(`sitemap/generate/error_email_identity`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SitemapGenerateErrorEmailTemplate = model.NewStr(`sitemap/generate/error_email_template`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SitemapGenerateFrequency = model.NewStr(`sitemap/generate/frequency`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SitemapGenerateTime = model.NewStr(`sitemap/generate/time`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SitemapLimitMaxLines = model.NewStr(`sitemap/limit/max_lines`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SitemapLimitMaxFileSize = model.NewStr(`sitemap/limit/max_file_size`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SitemapSearchEnginesSubmissionRobots = model.NewBool(`sitemap/search_engines/submission_robots`, model.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
