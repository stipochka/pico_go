#pragma once

#include <stdint.h>
#include <stddef.h>

#define PDU_BUFFER_SIZE 128
#define CRC8_POLYNOM 0x07
#define CRC8_INIT 0x00
#define TCP_PROTOCOl_CODE 6

// Ошибки picolib
typedef enum {
    PL_ERROR_NONE,                      //< Отсутствие ошибок
    PL_ERROR_NO_PACKAGE,                //< Отсутствие входящего пакета

    PL_ERROR_INVALID_SENSOR_NAME,       //< Неправильное имя сенсора
    PL_ERROR_INVALID_ARGUMENT,          //< Неправильный(-ые) аргументы
    PL_ERROR_INVALID_FUNCTION_CODE,     //< Неизвестный код функции

    PL_ERROR_INVALID_CRC,               //< Несовпадение CRC суммы
    PL_ERROR_REQUEST_TIMEOUT,           //< Таймаут запроса
    PL_ERROR_TRANSPORT,                 //< Ошибка транспорта
} pl_error_t;

// Функции picolib
typedef enum {
    PL_FUNC_HEARTBIT,                   //< Хартбит
    PL_FUNC_GET_ACTUAL_DATA,            //< Получить последнее значение с датчика
    PL_FUNC_GET_HISTORY_DATA,           //< Получить несколько последних значений с датчика
    PL_FUNC_GET_SENSOR_INFO,            //< Получить информацию о сенсоре
    PL_FUNC_GET_MCU_INFO,               //< Получить информацию о MCU
    PL_FUNC_SET_READING_PERIOD,         //< Установить период считывания датчика
} pl_func_t;

// Структура, хранящая данные PDU
#pragma pack(push, 1)
typedef struct {
    uint8_t  func_code;                 //< Код функции
    uint8_t  err_code;                  //< Код ошибки
    uint16_t args;                      //< Аргументы
    uint8_t  buffer[PDU_BUFFER_SIZE];   //< Буффер
} pl_pdu_t;

typedef struct {
    pl_pdu_t pdu;                       //< PDU
    uint8_t crc;                        //< Контрольная сумма
} pl_adu_t;
#pragma pack(pop)

typedef struct transport_config {
    pl_error_t (*pl_send_fn)(const uint8_t* data, size_t data_len);
    pl_error_t (*pl_read_fn)(uint8_t* data, size_t* data_len);
} pl_transport_config_t;

/// @brief Функция вычисления контрольной суммы CRC-8
/// @param data   данные
/// @param length длина данных
/// @return один байт чек-суммы
uint8_t calculate_crc8(const uint8_t *data, size_t length);

/// @brief Функция для вычисления 1-байтовой контрольной суммы TCP
/// @param data   данные
/// @param length длина данных
/// @param src_ip ip-адрес отправителя
/// @param dst_ip ip-адрес получателя
/// @return один байт чек-суммы
uint8_t tcp_checksum(const void *data, size_t length, uint32_t src_ip, uint32_t dst_ip);

/// @brief Запаковывает буфер в ADU
/// @param adu
/// @param func_code
/// @param err_code
/// @param args
/// @param data
/// @param data_length
void pl_pack(pl_adu_t *adu, uint8_t func_code, uint8_t err_code, uint16_t args, const uint8_t *data, size_t data_length);

/// @brief Распаковывает adu в буффер (конвертер буффер 2 pdu)
/// @param adu
/// @param data
/// @param data_length
void pl_unpack(pl_adu_t *adu, const uint8_t *data, size_t data_length);

/// @brief
/// @param data
/// @param data_length
/// @return Код ошибки
pl_error_t pl_send(const uint8_t* data, size_t data_length);

/// @brief
/// @param data
/// @param data_length
/// @return Код ошибки
pl_error_t pl_recv(uint8_t* data, size_t* data_length);

/// @brief Сеттер транспорта отправки сообщения
/// @param pl_write Указатель на пользовательскую функцию отправки пакетов
pl_error_t pl_set_transport_send(pl_error_t (*pl_write)(const uint8_t* data, size_t data_len));

/// @brief Сеттер транспорта получения сообщения
/// @param pl_read Указатель на пользовательскую функцию чтения пакетов
pl_error_t pl_set_transport_recv(pl_error_t (*pl_read)(uint8_t* data, size_t* data_len));

/// @brief
/// @param func_code
/// @param err_code
/// @param args
/// @return Код ошибки
pl_error_t pl_reader(uint8_t* func_code, uint8_t* err_code, uint16_t* args, uint8_t* data, size_t data_length);

/// @brief
/// @param transport_config Транспортный конфиг
/// @return Код ошибки
pl_error_t init_pico(pl_transport_config_t* transport_config);
