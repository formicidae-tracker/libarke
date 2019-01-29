#include <gtest/gtest.h>

#include <vector>
#include <cmath>

#include "arke-private-conversion.h"






TEST(Conversion,Tmp1075ToFloat) {
	struct TestData {
		float FloatValue;
		uint16_t BinaryValue;
	};

	//extracted data from the tmp1075 datasheet, with the last four bit removed (always equal to zero)
	std::vector<TestData> testData = {
		{127.9375,0x7FF},
		{100,0x640},
		{80,0x500},
		{75,0x4B0},
		{50,0x320},
		{25,0x190},
		{0.25,0x004},
		{0.0625,0x001},
		{0,0x000},
		{-0.0625,0xFFF},
		{-0.25,0xFFC},
		{-25,0xE70},
		{-50,0xCE0},
		{-128,0x800}
	};

	for (auto d : testData) {
		EXPECT_EQ(tmp1075_to_float(d.BinaryValue),d.FloatValue)
			<< "While converting 0x"<< std::hex << d.BinaryValue;
	}
}


TEST(Conversion,FloatToHumidity) {
	struct TestData {
		float FloatValue;
		uint16_t BinaryValue;
	};

	std::vector<TestData> testData = {
		{101.0,16382},
		{100.00,16382},
		{0.0,0},
		{-1.0,0},
		{50.0,16382/2}
	};

	for (auto d : testData) {
		EXPECT_EQ(float_to_humidity(d.FloatValue),d.BinaryValue)
			<< "While converting " << d.FloatValue << "% RH";
	}

}


TEST(Conversion,HumidityToFloat) {
	struct TestData {
		float FloatValue;
		uint16_t BinaryValue;
	};

	EXPECT_TRUE(std::isnan(humidity_to_float(16383)));
	EXPECT_TRUE(std::isnan(humidity_to_float(16384)));

	std::vector<TestData> testData = {
		{100.00,16382},
		{0.0,0},
		{50.0,16382/2}
	};

	for (auto d : testData) {
		EXPECT_EQ(humidity_to_float(d.BinaryValue),d.FloatValue)
			<< "While converting " << d.BinaryValue;
	}

}

TEST(Conversion,FloatToTemperatureHIH) {
	struct TestData {
		float FloatValue;
		uint16_t BinaryValue;
	};

	std::vector<TestData> testData = {
		{126.0,16382},
		{125.00,16382},
		{0.0,3971},
		{-40,0},
		{-41.0,0},
	};

	for (auto d : testData) {
		EXPECT_EQ(float_to_hih6030_temperature(d.FloatValue),d.BinaryValue)
			<< "While converting " << d.FloatValue << "°C";
	}

}


TEST(Conversion,HIHTemperatureToFloat) {
	struct TestData {
		float FloatValue;
		uint16_t BinaryValue;
	};

	EXPECT_TRUE(std::isnan(hih6030_temperature_to_float(16383)));
	EXPECT_TRUE(std::isnan(hih6030_temperature_to_float(16384)));

	std::vector<TestData> testData = {
		{125.00,16382},
		{0.0,3971},
		{-40.0,0}
	};

	for (auto d : testData) {
		EXPECT_TRUE(std::fabs(hih6030_temperature_to_float(d.BinaryValue)-d.FloatValue) < 0.005)
			<< "While converting " << d.BinaryValue;
	}

}
