package ebml

const (
	// Rudimentary EBML Elements:
	//   https://matroska-org.github.io/libebml/specs.html
	EBML                   = 0x1A45DFA3
	EBML_VERSION           = 0x4286
	EBML_READ_VERSION      = 0x42F7
	EBML_MAX_ID_LENGTH     = 0x42F2
	EBML_MAX_SIZE_LENGTH   = 0x42F3
	DOCTYPE                = 0x4282
	DOCTYPE_VERSION        = 0x4287
	DOTYPE_READ_VERSION    = 0x4285
	CRC32                  = 0xBF
	EBML_VOID              = 0xEC
	SIGNATURE_SLOT         = 0x1B538667
	SIGNATURE_ALGO         = 0x738A
	SIGNATURE_PUBLIC_KEY   = 0x7EA5
	SIGNATURE              = 0x7EB5
	SIGNATURE_ELEMENTS     = 0x7E5B
	SIGNATURE_ELEMENT_LIST = 0x7E7B
	SIGNED_ELEMENT         = 0x6532

	// Matroska EBML Formats:
	//   https://www.matroska.org/technical/specs/index.html
	EBML_SEGMENT                    = 0x18538067
	EBML_INFO                       = 0x1549A966
	EBML_TIMECODE_SCALE             = 0x2AD7B1
	EBML_TIMECODE_SCALE_NUMERATOR   = EBML_TIMECODE_SCALE
	EBML_TIMECODE_SCALE_DENOMINATOR = 0x2AD7B2
	EBML_CUE_TIME                   = 0xB3
	EBML_CUE_TRACK_POSITIONS        = 0xB7
	EBML_CUE_CLUSTER_POSITION       = 0xF1
	EBML_DURATION                   = 0x4489
	EBML_CUES                       = 0x1C53BB6B
	EBML_CUE_POINT                  = 0xBB
)
