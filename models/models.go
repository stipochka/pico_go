package models

/*
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
*/

type McuResponse struct {
	FuncCode uint8
	ErrCode  uint8
	Args     uint16
	Buffer   []byte
}
