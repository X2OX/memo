package participle

import (
	"path"
	"strings"

	"github.com/yanyiwu/gojieba"
)

func Init(folderName string) {
	gojieba.DICT_DIR = path.Join(path.Dir(folderName), "dict")
	gojieba.DICT_PATH = path.Join(gojieba.DICT_DIR, "jieba.dict.utf8")
	gojieba.HMM_PATH = path.Join(gojieba.DICT_DIR, "hmm_model.utf8")
	gojieba.USER_DICT_PATH = path.Join(gojieba.DICT_DIR, "user.dict.utf8")
	gojieba.IDF_PATH = path.Join(gojieba.DICT_DIR, "idf.utf8")
	gojieba.STOP_WORDS_PATH = path.Join(gojieba.DICT_DIR, "stop_words.utf8")
}

func Parse(s string) string {
	x := gojieba.NewJieba()
	defer x.Free()

	return strings.Join(x.CutForSearch(s, true), " ")
}
