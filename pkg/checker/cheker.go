package checker

import "os"

type Checker struct {
}

/*
void find_next_nonspace(int & c1, int & c2, FILE *& f1, FILE *& f2, int & ret){
    // Find the next non-space character or \n.
    while ((isspace(c1)) || (isspace(c2))){
        if (c1 != c2){
            if (c2 == EOF){
                do{
                    c1 = fgetc(f1);
                }while (isspace(c1));
                continue;
            }else if (c1 == EOF){
                do{
                    c2 = fgetc(f2);
                }while (isspace(c2));
                continue;
            }else if ((c1 == '\r' && c2 == '\n')){
                c1 = fgetc(f1);
            }else if ((c2 == '\r' && c1 == '\n')){
                c2 = fgetc(f2);
            }else{
                ret = JudgePE;
            }
        }
        if (isspace(c1)){
            c1 = fgetc(f1);
        }if (isspace(c2)){
            c2 = fgetc(f2);
        }
    }
}


/*
 * translated from ZOJ judger r367
 * http://code.google.com/p/zoj/source/browse/trunk/judge_client/client/text_checker.cc#25
 *
*/
//比较用户输出和标准数据
/*
int compare(const char *file1, const char *file2)
{
    int ret = JudgeAC;
    int c1,c2;
    FILE * f1, *f2 ;
    f1 = fopen(file1, "r");
    f2 = fopen(file2, "r");
    if (!f1 || !f2)
    {
        ret = JudgeRE;
    }
    else
        for (;;)
        {
            // Find the first non-space character at the beginning of line.
            // Blank lines are skipped.
            c1 = fgetc(f1);
            c2 = fgetc(f2);
            find_next_nonspace(c1, c2, f1, f2, ret);
            // Compare the current line.
            for (;;)
            {
                // Read until 2 files return a space or 0 together.
                while ((!isspace(c1) && c1) || (!isspace(c2) && c2))
                {
                    if (c1 == EOF && c2 == EOF)
                    {
                        goto end;
                    }
                    if (c1 == EOF || c2 == EOF)
                    {
                        break;
                    }
                    if (c1 != c2)
                    {
                        // Consecutive non-space characters should be all exactly the same
                        ret = JudgeWA;
                        goto end;
                    }
                    c1 = fgetc(f1);
                    c2 = fgetc(f2);
                }
                find_next_nonspace(c1, c2, f1, f2, ret);
                if (c1 == EOF && c2 == EOF)
                {
                    goto end;
                }
                if (c1 == EOF || c2 == EOF)
                {
                    ret = JudgeWA;
                    goto end;
                }

                if ((c1 == '\n' || !c1) && (c2 == '\n' || !c2))
                {
                    break;
                }
            }
        }
end:
    // if(ret==JudgeWA)make_diff_out(f1,f2,c1,c2,file1);
    if (f1)
        fclose(f1);
    if (f2)
        fclose(f2);
    return ret;
}
*/

func (ck *Check) Check(a, b *os.File) (int, error) {

}
