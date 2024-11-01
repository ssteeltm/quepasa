package models

import (
	"fmt"
	"path/filepath"
	"strings"

	library "github.com/nocodeleaks/quepasa/library"
	whatsapp "github.com/nocodeleaks/quepasa/whatsapp"
)

func SecureAndCustomizeAttach(attach *whatsapp.WhatsappAttachment) (extra []string) {
	if attach == nil {
		return
	}

	var contentMime string
	content := attach.GetContent()
	if content != nil {
		contentMime = library.GetMimeTypeFromContent(*content)
		extra = append(extra, fmt.Sprintf("[debug][SecureAndCustomizeAttach] detected mime type from content: %s", contentMime))
	}

	var requestExtension string
	if len(attach.FileName) > 0 {
		fileNameNormalized := strings.ToLower(attach.FileName) // important, because some bitches do capitalize filenames
		requestExtension = filepath.Ext(fileNameNormalized)
		extra = append(extra, fmt.Sprintf("[debug][SecureAndCustomizeAttach] detected extension: %s, from filename: %s", requestExtension, attach.FileName))
	} else if len(attach.Mimetype) > 0 {
		requestExtension, _ = library.TryGetExtensionFromMimeType(attach.Mimetype)
		extra = append(extra, fmt.Sprintf("[debug][SecureAndCustomizeAttach] detected extension from request mime type: %s", requestExtension))
	} else if len(contentMime) > 0 {
		requestExtension, _ = library.TryGetExtensionFromMimeType(contentMime)
		extra = append(extra, fmt.Sprintf("[debug][SecureAndCustomizeAttach] detected extension from content mime type: %s", requestExtension))
	}

	if len(contentMime) > 0 {

		if strings.HasPrefix(contentMime, "text/xml") && requestExtension == ".svg" {
			contentMime = "image/svg+xml"
		}

		if len(attach.Mimetype) == 0 {
			attach.Mimetype = contentMime
			extra = append(extra, fmt.Sprintf("[debug][SecureAndCustomizeAttach] updating empty mime type from content: %s", contentMime))
		}

		contentExtension, success := library.TryGetExtensionFromMimeType(contentMime)
		if success {
			extra = append(extra, fmt.Sprintf("[debug][SecureAndCustomizeAttach] content extension: %s", contentExtension))

			// if was passed a filename without extension
			if len(requestExtension) == 0 && len(attach.FileName) > 0 {
				extra = append(extra, fmt.Sprintf("[info][SecureAndCustomizeAttach] missing extension for attachment (%s), using from content: %s :: content mime: %s", attach.FileName, contentExtension, contentMime))

				attach.Mimetype = contentMime
				attach.FileName += contentExtension
			} else {
				// validating mime information
				if !IsValidExtensionFor(requestExtension, contentExtension) {
					// invalid attachment
					extra = append(extra, fmt.Sprintf("[warn][SecureAndCustomizeAttach] invalid extension for attachment, request extension: %s (%s) != content extension: %s :: content mime: %s, revalidating for security", requestExtension, attach.FileName, contentExtension, contentMime))
					attach.Mimetype = contentMime
					attach.FileName = whatsapp.InvalidFilePrefix + library.GenerateFileNameFromMimeType(contentMime)
				}
			}
		}
	}

	// set compatible audios to be sent as ptt
	ForceCompatiblePTT := ENV.UseCompatibleMIMEsAsAudio()
	if ForceCompatiblePTT && !attach.IsValidAudio() && IsCompatibleWithPTT(attach.Mimetype) {
		extra = append(extra, fmt.Sprintf("[info][SecureAndCustomizeAttach] setting that it should be sent as ptt, regards its incompatible mime type: %s", attach.Mimetype))
		attach.SetPTTCompatible(true)
	}

	// Defining a filename if not found before
	if len(attach.FileName) == 0 {
		attach.FileName = library.GenerateFileNameFromMimeType(attach.Mimetype)
		extra = append(extra, fmt.Sprintf("[debug][SecureAndCustomizeAttach] empty file name, generating a new one based on mime type: %s, file name: %s", attach.Mimetype, attach.FileName))
	}

	extra = append(extra, fmt.Sprintf("[debug][SecureAndCustomizeAttach] resolved mime type: %s, filename: %s", attach.Mimetype, attach.FileName))
	return
}

func IsValidExtensionFor(request string, content string) bool {
	switch {
	case
		request == ".jpg" && content == ".jpeg", // used for correct old windows 3 characters extensions
		request == ".csv" && content == ".txt",
		request == ".json" && content == ".txt",
		request == ".sql" && content == ".txt",
		request == ".oga" && content == ".webm",
		request == ".oga" && content == ".ogx",
		request == ".opus" && content == ".ogx",
		request == ".ovpn" && content == ".txt",
		request == ".svg" && content == ".xml":
		return true
	}

	return request == content
}

func IsCompatibleWithPTT(mime string) bool {
	// switch for basic mime type, ignoring suffix
	mimeOnly := strings.Split(mime, ";")[0]

	for _, item := range whatsapp.WhatsappMIMEAudioPTTCompatible {
		if item == mimeOnly {
			return true
		}
	}

	return false
}
