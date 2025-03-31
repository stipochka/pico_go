#pragma once

#include "picolib.h"

/// @brief
/// @param
/// @return Код ошибки
pl_error_t heartbit_request(uint8_t* data, uint16_t *data_size);

/// @brief
/// @param
/// @return Код ошибки
pl_error_t get_actual_data_request(const uint8_t* sensor_name, uint8_t* data, uint16_t *data_size);

/// @brief
/// @param
/// @return Код ошибки
pl_error_t get_history_data_request(const uint8_t* sensor_name, uint16_t num, uint8_t* data, uint16_t *data_size);

/// @brief
/// @param
/// @return Код ошибки
pl_error_t get_sensor_info_request(const uint8_t* sensor_name, uint8_t* data, uint16_t *data_size);

/// @brief
/// @param
/// @return Код ошибки
pl_error_t get_mcu_info_request(uint8_t* data, uint16_t *data_size);

/// @brief
/// @param
/// @return Код ошибки
pl_error_t set_reading_period_request(const uint8_t* sensor_name, uint16_t delay, uint8_t* data, uint16_t *data_size);
