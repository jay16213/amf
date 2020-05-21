package callback

import (
	"context"
	"free5gc/lib/openapi/Namf_Communication"
	"free5gc/lib/openapi/models"
	amf_context "free5gc/src/amf/context"
	"free5gc/src/amf/logger"
	"reflect"
)

func SendAmfStatusChangeNotify(amfStatus string, guamiList []models.Guami) {
	amfSelf := amf_context.AMF_Self()

	for _, amfStatusSubscriptions := range amfSelf.AMFStatusSubscriptions {

		configuration := Namf_Communication.NewConfiguration()
		client := Namf_Communication.NewAPIClient(configuration)
		amfStatusNotification := models.AmfStatusChangeNotification{}
		var amfStatusInfo = models.AmfStatusInfo{}

		for _, guami := range guamiList {
			for _, subGumi := range amfStatusSubscriptions.GuamiList {
				if reflect.DeepEqual(guami, subGumi) {
					//AMF status is available
					amfStatusInfo.GuamiList = append(amfStatusInfo.GuamiList, guami)

				}
			}
		}

		amfStatusInfo = models.AmfStatusInfo{
			StatusChange:     (models.StatusChange)(amfStatus),
			TargetAmfRemoval: "",
			TargetAmfFailure: "",
		}

		amfStatusNotification.AmfStatusInfoList = append(amfStatusNotification.AmfStatusInfoList, amfStatusInfo)
		uri := amfStatusSubscriptions.AmfStatusUri

		logger.ProducerLog.Infof("[AMF] Send Amf Status Change Notify to %s", uri)
		httpResponse, err := client.AmfStatusChangeCallbackDocumentApiServiceCallbackDocumentApi.AmfStatusChangeNotify(context.Background(), uri, amfStatusNotification)
		if err != nil {
			if httpResponse == nil {
				HttpLog.Errorln(err.Error())
			} else if err.Error() != httpResponse.Status {
				HttpLog.Errorln(err.Error())
			}
		}
	}
}
