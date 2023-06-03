from datetime import datetime
import pytz
import csv
import os

# Compare timezone data
def compareTimezones(inputDatetime, tzPST, tzEST, rowData):
	try:
		# Get timezone to compare
		tzToCheck = pytz.timezone(str(rowData[0]))
		print(f'Timezone to Compare: {tzToCheck}')
		tzCompareLocalized = tzToCheck.localize(inputDatetime)
		print(f'Compared Timezone Time: {tzCompareLocalized}')

		# Compare against PST
		comparePST = abs((tzPST - tzCompareLocalized).total_seconds()/3600)
		print(f'Variation to PST Time: {comparePST}')

		# Compare against EST
		compareEST = abs((tzEST - tzCompareLocalized).total_seconds()/3600)
		print(f'Variation to EST Time: {compareEST}')

		# Compare the PST vs EST
		comparison = "None"
		if comparePST > compareEST:
			comparison = "Closer to EST"
		elif comparePST < compareEST:
			comparison = "Closer to PST"
		else:
			comparison = "No difference can be used either in EST or PST"

		# Add to the array
		csv_processed_data = [rowData[0], comparePST, compareEST, comparison]
		return csv_processed_data
	
	except Exception as err:
		print(f'Error Occurred: {err}')
	
# Write processed data to a new CSV file
def writeProcessedData(data):
	try:
		with open(f'{os.getcwd()}/python/processed_data/timezone_data.csv', 'w', newline='') as csv_file:
			csv_writer = csv.writer(csv_file, delimiter=',')
			# Write the header row
			csv_writer.writerow(['Timezone', 'Variation to PST', 'Variation to EST', 'Comparison'])
			csv_writer.writerows(data)
	except Exception as err:
		print(f'Error Occurred: {err}')


def main():

	# Set primary timezones
	pstTZ = pytz.timezone('US/Pacific')
	estTZ = pytz.timezone('US/Eastern')
	inputDatetime = datetime.now()

	# Get the time in the required timezones
	print(f'Current Time: {inputDatetime}')
	pstTimezone = pstTZ.localize(inputDatetime)
	print(f'PST Time: {pstTimezone}')
	estTimezone = estTZ.localize(inputDatetime)
	print(f'EST Time: {estTimezone}')

	try:
		timezoneData = []

		# Read the source CSV file
		with open(f'{os.getcwd()}/python/data/timedata.csv') as csv_file:
			csv_reader = csv.reader(csv_file, delimiter=',')

			# Compare the timezones for each row in the CSV file
			# Add the processed data to the array
			for rowData in csv_reader:
				timezoneData.append(compareTimezones(inputDatetime, pstTimezone, estTimezone, rowData))

		# Call the function to write the processed data to a new CSV file
		writeProcessedData(timezoneData)

	except Exception as err:
		print(f'Error Occurred: {err}')
  
  
if __name__=="__main__":
    main()