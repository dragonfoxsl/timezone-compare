from datetime import datetime
import pytz
import csv

def compareTimezones(inputDatetime, tzPST, tzEST, rowData):
	tzToCheck = pytz.timezone(str(rowData[0]))
	print(f'Timezone to Compare: {tzToCheck}')
	tzCompareLocalized = tzToCheck.localize(inputDatetime)
	print(f'Timezone Time: {tzCompareLocalized}')

	comparePST = abs((tzPST - tzCompareLocalized).total_seconds()/3600)
	print(f'Variation to PST Time: {comparePST}')
	compareEST = abs((tzEST - tzCompareLocalized).total_seconds()/3600)
	print(f'Variation to EST Time: {compareEST}')

	comparison = "None"
	if comparePST > compareEST:
		comparison = "Closer to EST"
	elif comparePST < compareEST:
		comparison = "Closer to PST"
	else:
		comparison = "No difference can be used either in EST or PST"

	csv_processed_data = [rowData[0], comparePST, compareEST, comparison]
	return csv_processed_data 

def writeProcessedData(data):
    with open('processed_data/timezone_data.csv', 'w', newline='') as csv_file:
     csv_writer = csv.writer(csv_file, delimiter=',')
     csv_writer.writerow(['Timezone', 'Variation to PST', 'Variation to EST', 'Comparison'])
     csv_writer.writerows(data)

pstTZ = pytz.timezone('US/Pacific')
estTZ = pytz.timezone('US/Eastern')
inputDatetime = datetime.now()
print(f'Current Time: {inputDatetime}')
pstTimezone = pstTZ.localize(inputDatetime)
print(f'PST Time: {pstTimezone}')
estTimezone = estTZ.localize(inputDatetime)
print(f'EST Time: {estTimezone}')

timezoneData = []
with open('data/timedata.csv') as csv_file:
	csv_reader = csv.reader(csv_file, delimiter=',')
	
	for rowData in csv_reader:
		timezoneData.append(compareTimezones(inputDatetime, pstTimezone, estTimezone, rowData))

writeProcessedData(timezoneData)